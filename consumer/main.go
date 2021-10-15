package main

import (
	"consmer/config"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strings"
)

const updateLinks = "http://%s/api/v1/links/%d"

type ReqLink struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func main() {
	cfg := config.Init()

	conn, err := amqp.Dial(cfg.QueueAdress)
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
			status := GetUrlStatus(reqLink.URL)

			err, status = SendUpdateLinkRequest(reqLink.ID, status, cfg)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("URL: %s STATUS_CODE: %d", reqLink.URL, status)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
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

func SendUpdateLinkRequest(id int, status int, cfg *config.Config) (error, int) {
	adress := fmt.Sprintf(updateLinks, cfg.ServerAdress, id)
	client := resty.New()
	_, err := client.R().SetBody(struct {
		Status int `json:"status_code"`
	}{status},
	).Put(adress)
	if err != nil {
		return err, 0
	}
	return nil, status
}
