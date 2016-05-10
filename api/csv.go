package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Payload is used to organize the incoming json data
type Payload struct {
	Filename  string `json:"filename"`
	Path      string `json:"path"`
	Extension string `json:"extension"`
	Key       string `json:"key"`
	Secret    string `json:"secret"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
}

func csvAsync(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(500)
	} else {
		file := Payload{}
		json.Unmarshal([]byte(b), &file)

		f := &AwsInfo{
			AccessKey: file.Key,
			SecretKey: file.Secret,
			Region:    file.Region,
			Bucket:    file.Bucket,
		}

		localFile := "/tmp/" + file.Filename + "." + file.Extension
		c := awsConnect(f)
		awsDownload(c, string(file.Path+"/"+file.Filename+"."+file.Extension), localFile, fileQueue)
	}
}

// createRowToJson receives a column slice and a values slice, then matches each column-value and outputs json
func createRowToJson(columns []string, values []string) ([]byte, error) {
	// create a map using exactly the values' length
	mp := make(map[string]string, len(values))

	for i, v := range values {
		mp[columns[i]] = v
	}

	return json.Marshal(mp)
}

func csvAsyncProcessor(ch chan string) {
	for m := range ch {
		f, err := os.Open(m)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		r := csv.NewReader(f)
		r.FieldsPerRecord = -1
		rows, err := r.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Processing %s with %d rows", f, len(rows))

		// we expect the columns in the very first row
		columns := rows[0]

		// subslice rows in order to skip the first row (the columns)
		for _, row := range rows[1:len(rows)] {
			json, err := createRowToJson(columns, row)

			if err != nil {
				log.Println(err)
			} else {
				err := publish(rmqConn, os.Getenv("MESSAGEQUEUESERVER_EXCHANGE"), os.Getenv("MESSAGEQUEUESERVER_QUEUE"), string(json), true)
				if err != nil {
					log.Println(err)
				}

			}
		}

	}
}
