package main

import (
	"consumer/config"
	"consumer/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const (
	updateLinks    = "http://%s/api/v1/links/%d"
	foundRedisText = "found URL: %s STATUS_CODE: %d IN REDIS"
	answerText     = "URL: %s STATUS_CODE: %d"
	waitingReqMsg  = " [*] Waiting for messages."
)

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

	if err := ch.ExchangeDeclare("remind", "direct", false, true, false, false, nil); err != nil {
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
		"remind",      // sourceExchange
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
	dbService, err := db.NewPgxCon(ctx, cfg)

	go func() {
		for d := range msgs {
			var remind db.Remind
			err := json.Unmarshal(d.Body, &remind)
			if err != nil {
				log.Println(err)
				continue
			}

			if err := dbService.AddRemind(remind); err != nil {
				log.Println(err)
				continue
			}

			fmt.Printf("Напоминание успешно добавлено из очереди")
		}
	}()

	log.Printf(waitingReqMsg)
	<-forever
}
