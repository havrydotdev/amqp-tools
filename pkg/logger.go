package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/gavrylenkoIvan/amqp-tools/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Tools struct {
	rabbit *amqp.Connection
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type EmailData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Url      string `json:"url"`
	Delete   string `json:"delete"`
}

// Returns new logger. Gets rabbitmq url as argument
func New(url string) (*Tools, error) {
	conn, err := connect(url)
	if err != nil {
		return nil, err
	}

	return &Tools{
		rabbit: conn,
	}, nil
}

func connect(url string) (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	// don`t continue until rabbit is ready
	for {
		c, err := amqp.Dial(url)
		if err != nil {
			fmt.Println("Failed to connect to rabbitmq: ", err)
			counts++
		} else {
			log.Println("Connected to rabbitmq")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("Tried 5 times to connect to rabbitmq, giving up")
			return nil, errors.New("failed to connect to rabbitmq")
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off for", backoff.Seconds())
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}

func (t *Tools) WriteLog(name, msg string) error {
	return t.sendMessage(
		LogPayload{name, msg},
		"logs_topic", "log",
	)
}

func (t *Tools) SendEmail(data EmailData) error {
	return t.sendMessage(data, "logs_topic", "email.AUTH")
}

func (t *Tools) sendMessage(data any, exchange, severity string) error {
	emitter, err := event.NewEventEmitter(t.rabbit)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(data)
	err = emitter.Push(string(j), severity, exchange)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tools) WriteErrorLog(name, msg string) {
	err := t.WriteLog(name, msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(name, msg)
}
