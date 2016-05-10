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

func consume(conn *amqp.Connection, exchange, queue, consumerTag string) error {

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: Could not declare channel %q", err)
		return err
	} else {

		log.Printf("Got Channel")
		defer ch.Close()

		err := ch.ExchangeDeclare(
			exchange, // name
			"topic",  // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // noWait
			nil,      // arguments
		)
		if err != nil {
			log.Printf("ERROR: Could not declare Exchange [%s] %q", exchange, err)
			return nil
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
				log.Printf("ERROR: Could not declare queue [%s] %q", queue, err)
				return err
			} else {

				err = ch.QueueBind(
					q.Name,   // queue name
					queue,    // routing key @TODO
					exchange, // exchange
					false,
					nil,
				)
				if err != nil {
					log.Printf("ERROR: Could not bind [%s] queue to [%s] exhange %q", queue, exchange, err)
					return err
				} else {

					log.Printf("Queue %s declared", queue)
					err = ch.Qos(
						1,     // prefetch count
						0,     // prefetch size
						false, // global
					)
					if err != nil {
						log.Printf("ERROR: %q", err)
						return err
					}
					msgs, err := ch.Consume(
						q.Name,      // queue
						consumerTag, // consumer
						false,       // auto-ack
						false,       // exclusive
						false,       // no-local
						false,       // no-wait
						nil,         // args
					)
					if err != nil {
						log.Printf("ERROR: Could not consume messages on [%s] queue %q", queue, err)
						return err
					}

					forever := make(chan bool)

					go func() {
						for d := range msgs {
							d.Ack(false)

							// Dropped in here for now but need to push d.Body to a channel
							err = indexData(string(d.Body))
						}
					}()
					log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
					<-forever
					return nil
				}

			}

		}

	}
	return nil
}
