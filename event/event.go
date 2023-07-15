package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	exchanges = []string{
		"logs_topic",
		"emails_topic",
	}
)

func declareExchange(ch *amqp.Channel) error {
	for _, ex := range exchanges {
		err := ch.ExchangeDeclare(
			ex,      // name
			"topic", // type
			true,    // durable
			false,   // auto-deleted
			false,   // internal
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			return err
		}
	}

	return nil
}
