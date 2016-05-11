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

// Not a fan of this func, I want to write a lib to handle all rmq actions...
func publish(connection *amqp.Connection, exchange, queue, body string, reliable bool) error {

	ch, err := connection.Channel()
	if err != nil {
		log.Printf("ERROR: Could not create channel - %s", err)
		return err
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
			log.Printf("ERROR: Could not declare exchange '%s' - %s", exchange, err)
			return err
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
				log.Printf("ERROR: Could not declare queue '%s' - %s", queue, err)
				return err
			} else {

				err = ch.QueueBind(
					q.Name, // queue name
					queue,
					exchange, // exchange
					false,
					nil,
				)
				if err != nil {
					log.Printf("ERROR: Could not bind queue '%s' to exchange '%s' using '%s'", queue, exchange, queue, err)
					return err
				} else {

					// Reliable publisher confirms require confirm.select support from to connection.
					if reliable {
						if err := ch.Confirm(false); err != nil {
							log.Printf("ERROR: Could not confirm - %s", err)
							return err
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
						log.Printf("ERROR: Could not publish message %s", err)
						return err
					} else {
						ch.Close()
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
