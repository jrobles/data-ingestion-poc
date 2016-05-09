package main

import (
	"fmt"
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

func publish(connection *amqp.Connection, exchange, queue, routingKey, body string, reliable bool) error {

	ch, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	} else {

		if err := ch.ExchangeDeclare(
			exchange, // name
			"topic",  // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // noWait
			nil,      // arguments
		); err != nil {
			return fmt.Errorf("Exchange Declare: %s", err)
		} else {

			q, err := ch.QueueDeclare(
				queue, // name
				true,  // durable
				false, // delete when unused
				false, // exclusive
				false, // no-wait
				nil,   // arguments
			)
			if err != nil {
				return fmt.Errorf("Queue Declare: %s", err)
			} else {

				err = ch.QueueBind(
					q.Name, // queue name
					routingKey,
					exchange, // exchange
					false,
					nil,
				)
				if err != nil {
					log.Printf("ERROR: Could not bind [%s] queue to [%s] exhange %q", q.Name, exchange, err)
					return err
				} else {

					// Reliable publisher confirms require confirm.select support from the
					// connection.
					if reliable {
						if err := ch.Confirm(false); err != nil {
							return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
						}

						confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

						defer confirmOne(confirms)
					}

					if err = ch.Publish(
						exchange, // publish to an exchange
						q.Name,   // routing to 0 or more queues
						false,    // mandatory
						false,    // immediate
						amqp.Publishing{
							Headers:         amqp.Table{},
							ContentType:     "text/plain",
							ContentEncoding: "",
							Body:            []byte(body),
							DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
							Priority:        0,              // 0-9
							// a bunch of application/implementation-specific fields
						},
					); err != nil {
						return fmt.Errorf("Exchange Publish: %s", err)
					}

				}

			}
		}

	}

	return nil
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")
	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
