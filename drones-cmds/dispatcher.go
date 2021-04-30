package dronescmds

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type amqpDispatcher interface {
	Dispatch(queue string, event interface{}) error
}

type rabbitMQDispatcher struct {
	URL string
}

type rabbitMqConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func newRabbitMQDispatcher(queue string, url string) (dispatcher *rabbitMQDispatcher) {
	return &rabbitMQDispatcher{
		URL: url,
	}
}

func (d *rabbitMQDispatcher) Dispatch(queue string, event interface{}) (err error) {
	// establish rabbitmq connection
	conn, err := d.connect()
	defer conn.close()
	// send message
	msg, err := json.Marshal(&event)
	if err != nil {
		log.Println("Failed to encode message", err)
		return err
	}
	// create queue (indempotent: creates only if it doesn't exist)
	_, err = conn.Channel.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println("Failed to create the queue", err)
		return err
	}
	err = conn.Channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		log.Println("Failed to publish a message", err)
		return err
	}
	return
}

func (d *rabbitMQDispatcher) connect() (rabbitMqConnection, error) {
	// makes socket connection
	conn, err := amqp.Dial(d.URL)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ", err)
		return rabbitMqConnection{}, err
	}
	// create channel for messaging
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel", err)
		return rabbitMqConnection{}, err
	}

	return rabbitMqConnection{conn, ch}, nil
}

func (rc *rabbitMqConnection) close() {
	rc.Connection.Close()
	rc.Channel.Close()
}
