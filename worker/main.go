package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

var (
	rmqConn   *amqp.Connection = nil
	rmqServer                  = "amqp://guest:guest@" + os.Getenv("MESSAGEQUEUESERVER_PORT_5672_TCP_ADDR") + ":" + os.Getenv("MESSAGEQUEUESERVER_PORT_5672_TCP_PORT")
)

func init() {
	rmqConn = connect(rmqServer)
}

func main() {
	err := consume(rmqConn, os.Getenv("MESSAGEQUEUESERVER_EXCHANGE"), os.Getenv("MESSAGEQUEUESERVER_QUEUE"), "W1")
	if err != nil {
		log.Printf("ERROR: %q", err)
	}
}
