package main

import (
	"encoding/json"
	"github.com/mattbaird/elastigo/lib"
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

func indexData(body string) error {

	c := elastigo.NewConn()
	c.RequestTracer = func(method, url, body string) {
		log.Printf("Requesting %s %s", method, url)
		log.Printf("Request body: %s", body)
	}
	c.Domain = os.Getenv("ELASTICSEARCHSERVER_PORT_9200_TCP_ADDR")

	var objmap map[string]*json.RawMessage
	var pid string
	b := []byte(body)
	err := json.Unmarshal(b, &objmap)
	if err != nil {
		log.Printf("ERROR: Could not unmarshal message - %s", err)
	}

	// Can't hard-code the unique identifier, works for POC but needs to be fixed
	err = json.Unmarshal(*objmap["Item Number"], &pid)
	if err != nil {
		log.Printf("ERROR: Coulnd not unarhal Item Number -  %s", err)
	}
	_, err = c.Index("spacely_sprockets", "products", pid, nil, body)
	if err != nil {
		log.Printf("ERROR: Could not index into elasticsearch %s", err)
	}

	return nil
}
