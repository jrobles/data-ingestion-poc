package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/v1/csv", csvImport).Methods("POST")
	log.Fatal(http.ListenAndServe(":1337", router))
}
