package queue

import (
	"backend/internal/service/db"
	"encoding/json"
	"github.com/streadway/amqp"
)

func (q *RabbitCon) SetLinkStatus(id int, url string) error {
	ch, err := q.rbConn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}

	if err := ch.ExchangeDeclare("linker", "direct", false, true, false, false, nil); err != nil {
		return err
	}

	// We create a Queue to send the message to.
	//que, err := ch.QueueDeclare(
	//	"link-status", // name
	//	false,         // durable
	//	false,         // delete when unused
	//	false,         // exclusive
	//	false,         // no-wait
	//	nil,           // arguments
	//)
	//if err != nil {
	//	return err
	//}

	// We set the payload for the message.
	link := db.Link{ID: id, URL: url}

	body, _ := json.Marshal(&link)
	err = ch.Publish(
		"linker", // exchange
		"123",    // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
