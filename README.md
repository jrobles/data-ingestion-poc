# data-ingestion-poc
POC for a data ingestion microservice using Go, Elasticsearch, and Rabbitmq. The concept is: large feeds are imported via an API written in go which concurrently distributes the messages to N workers (Go) via rabbitmq. Each worker processes the record into elasticsearch.

Run the app via docker-compose
```
./bin/rebuild
```

Post to API
```
curl -X POST -H "Content-Type: application/json" -d '{
    "filename": "your_filename",
    "path": "/file/path/on/s3",
    "extension": "csv",
    "bucket": "your-bucket-name",
    "region": "your-region",
    "key": "YOUR_S3_KEY",
    "secret": "YOUR_S3_SECRET"
}' "http://{YOUR_DOCKER_MACHINE_IP}:1337/api/v1/csv"
```

Assuming the Item Number column had a value of 44870, the data can be obtained from Elasticsearch via:
```
curl -X GET "http://{YOUR_DOCKER_MACHINE_IP}:19200/{index}/{type}/_search?q=Item\ Number:440870"
```
yields:
```
{
  "took": 17,
  "timed_out": false,
  "_shards": {
    "total": 5,
    "successful": 5,
    "failed": 0
  },
  "hits": {
    "total": 1,
    "max_score": 8.711549,
    "hits": [
      {
        "_index": "spacely_sprockets",
        "_type": "products",
        "_id": "399088",
        "_score": 8.711549,
        "_source": {
          "Active Location": "3-A4-076-A-2",
          "Active Lock Code": "",
          "Active Units": "158",
          "Carton Units": "2",
          "Dept": "16",
          "Item Description": "TEAPOT ELEPHANT WHITE THING STYFF",
          "Item Number": "399088",
          "Receipt Date": "20140304",
          "Receipt ETA": "20140702",
          "Reserve Units": "2220",
          "SDC ECom Units": "0"
        }
      }
    ]
  }
}
```
