package main

import (
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

var (
	rmqConn *amqp.Connection = nil
)

func init() {
	rmqConn = connect("aadaadadsasasdaa.adasdadsa.hry")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/v1/csv", csvImport).Methods("POST")
	log.Fatal(http.ListenAndServe(":1337", router))
}
