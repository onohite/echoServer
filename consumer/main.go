package main

import (
	"consmer/config"
	"consmer/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	updateLinks = "http://%s/api/v1/links/%d"
	TTL         = time.Hour * 2
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
			cacheStatus, err := CheckCacheStatus(cfg, ctx, reqLink.URL)
			if err != nil {
				cacheStatus = GetUrlStatus(reqLink.URL)
				err := AddCacheStatus(cfg, ctx, reqLink.URL, cacheStatus)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			err = SendUpdateLinkRequest(reqLink.ID, cacheStatus, cfg)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("URL: %s STATUS_CODE: %d", reqLink.URL, cacheStatus)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func CheckCacheStatus(cfg *config.Config, ctx context.Context, uri string) (int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var cache redis.CacheDB
	cache.InitCache(cfg)
	result, err := cache.CacheConn.Get(ctx, uri).Result()
	if err != nil {
		return 0, err
	}
	status, _ := strconv.Atoi(result)
	return status, nil
}

func AddCacheStatus(cfg *config.Config, ctx context.Context, uri string, status int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var cache redis.CacheDB
	cache.InitCache(cfg)
	err := cache.CacheConn.Set(ctx, uri, status, TTL).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetUrlStatus(url string) int {
	if !strings.Contains("url", "http") && !strings.Contains("url", "https") {
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
