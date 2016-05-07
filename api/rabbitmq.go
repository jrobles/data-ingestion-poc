package main

import (
	"github.com/streadway/amqp"
	"log"
)

func connect(amqpURI string) (ch *amqp.Connection) {
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("ERROR: Could not connect to RabbitMQ server - %q", err)
		return nil
	} else {
		log.Printf("Connected to RabbitMQ server - %s", amqpURI)
		return connection
	}
}
