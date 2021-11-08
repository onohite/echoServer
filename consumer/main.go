package main

import (
	"consumer/config"
	"consumer/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strings"
)

const (
	updateLinks    = "http://%s/api/v1/links/%d"
	foundRedisText = "found URL: %s STATUS_CODE: %d IN REDIS"
	answerText     = "URL: %s STATUS_CODE: %d"
	waitingReqMsg  = " [*] Waiting for messages."
)

type ReqLink struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func main() {
	cfg := config.Init()

	conn, err := amqp.Dial(cfg.QueueAdress)
	ctx := context.Background()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare("linker", "direct", false, true, false, false, nil); err != nil {
		log.Fatalln(err)
	}

	q, err := ch.QueueDeclare(
		"link-status", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err = ch.QueueBind(
		"link-status", // name of the queue
		"123",         // bindingKey
		"linker",      // sourceExchange
		false,         // noWait
		nil,           // arguments
	); err != nil {
		log.Fatalf("Queue Bind: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalln(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var reqLink ReqLink
			err := json.Unmarshal(d.Body, &reqLink)
			if err != nil {
				log.Println(err)
				continue
			}
			redisService := redis.CacheDB{}
			redisService.InitCache(cfg)
			cacheStatus, err := redisService.CheckCacheStatus(cfg, ctx, reqLink.URL)
			if err != nil {
				log.Println(err)
				cacheStatus = GetUrlStatus(reqLink.URL)
				err := redisService.AddCacheStatus(cfg, ctx, reqLink.URL, cacheStatus)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				log.Printf(foundRedisText, reqLink.URL, cacheStatus)
			}

			err = SendUpdateLinkRequest(reqLink.ID, cacheStatus, cfg)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf(answerText, reqLink.URL, cacheStatus)
		}
	}()

	log.Printf(waitingReqMsg)
	<-forever
}

func GetUrlStatus(url string) int {
	if !strings.Contains(url, "http") || !strings.Contains(url, "https") {
		url = fmt.Sprintf("http://%s", url)
	}
	cl := resty.New()
	resp, err := cl.R().Get(url)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	return resp.StatusCode()
}

func SendUpdateLinkRequest(id int, status int, cfg *config.Config) error {
	adress := fmt.Sprintf(updateLinks, cfg.ServerAdress, id)
	client := resty.New()
	_, err := client.R().SetBody(struct {
		Status int `json:"status_code"`
	}{status},
	).Put(adress)
	if err != nil {
		return err
	}
	return nil
}
