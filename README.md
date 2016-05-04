# data-ingestion-poc
POC for a data ingestion microservice using Go, Elasticsearch, and Rabbitmq. The concept is: large feeds are imported via an API written in go which concurrently distributes the messages to N workers (Go) via rabbitmq. Each worker processes the record into elastic.
