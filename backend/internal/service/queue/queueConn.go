package queue

import (
	"context"
	"echoTest/internal/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const waitTimeout = 10 * time.Second

type RabbitCon struct {
	rbConn    *amqp.Connection
	rbConnCtx context.Context
}

func NewRabbitCon(ctx context.Context, config *config.Config) (*RabbitCon, error) {
	// Контекст ограниченный по времени ожидания
	instance := &RabbitCon{}

	instance.rbConnCtx = ctx

	var err error
	var count = 0
	for {
		if count < 6 {
			count++
		}
		err = instance.reconnect(config.QueueAdress)
		if err != nil {
			log.Printf("connection was lost. Error: %s. Wait %d sec.", err, count*5)
		} else {
			break
		}
		log.Println("Try to reconnect...")
		time.Sleep(time.Duration(count*5) * time.Second)
	}
	return instance, nil
}

func (db *RabbitCon) Close() error {
	db.rbConn.Close()
	return nil
}

func (db *RabbitCon) reconnect(address string) error {
	_, cancel := context.WithTimeout(db.rbConnCtx, waitTimeout)
	defer cancel()
	conn, err := amqp.Dial(address)
	if err != nil {
		return fmt.Errorf("unable to connection to database: %v", err)
	}

	db.rbConn = conn
	return nil
}
