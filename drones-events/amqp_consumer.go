package dronesevents

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type amqpConsumer struct {
	AmqpURL   string
	QueueName string
}

type rabbitMqConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (ac *amqpConsumer) Dequeue() (<-chan interface{}, error) {
	// establish rabbitmq connection
	conn, err := ac.Connect()
	if err != nil {
		log.Printf("Error connecting to the queue: %s", err)
		return nil, err
	}
	// create queue (indempotent: creates only if it doesn't exist)
	queue, err := conn.Channel.QueueDeclare(
		ac.QueueName, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Failed to create the queue: %s", err)
		return nil, err
	}
	events, err := conn.Channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("Error consuming messages from queue: %s, error:%s", ac.QueueName, err)
		return nil, err
	}
	// channel type conversion
	// pipes events from amqp channel to out channel for event-processing
	out := make(chan interface{})
	go func() {
		for evt := range events {
			var droneEvent interface{}
			err := json.Unmarshal(evt.Body, &droneEvent)
			if err != nil {
				log.Fatalf("Error unmarshaling event to JSON, event: %v, error: %s", evt.Body, err)
				continue
			}
			log.Printf("Consumed drone event: %s", droneEvent)
			out <- droneEvent
		}
		close(out)
	}()
	return out, nil
}

func (ac *amqpConsumer) Connect() (rabbitMqConnection, error) {
	// makes socket connection
	conn, err := amqp.Dial(ac.AmqpURL)
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
	log.Println("Connected to RabbitMQ")
	return rabbitMqConnection{conn, ch}, nil
}
