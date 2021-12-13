package queue

import (
	"backend/internal/service/db"
	"encoding/json"
	"github.com/streadway/amqp"
)

func (q *RabbitCon) SetRemind(rem db.Remind) error {
	ch, err := q.rbConn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}

	if err := ch.ExchangeDeclare("remind", "direct", false, true, false, false, nil); err != nil {
		return err
	}

	// We set the payload for the message.

	body, _ := json.Marshal(&rem)
	err = ch.Publish(
		"remind", // exchange
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
