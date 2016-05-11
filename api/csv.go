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

const (
	createRowToJSONErrorMismatchLength = "ERROR: Column's length must be, at least, equal to values' length"
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
		log.Printf("ERROR: Could not read POST data - %s", err)
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

// createRowToJSON receives a column slice and a values slice, then matches each column-value and outputs json
func createRowToJSON(columns []string, values []string) ([]byte, error) {

	if len(columns) < len(values) {
		return nil, fmt.Errorf(createRowToJSONErrorMismatchLength)
	}

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
			log.Printf("ERROR: Could not open %s file - %s", m, err)
		}
		defer f.Close()

		r := csv.NewReader(f)
		r.FieldsPerRecord = -1
		rows, err := r.ReadAll()
		if err != nil {
			log.Printf("ERROR: Could not read csv rows - %s", err)
			os.Exit(1)
		}

		// we expect the columns in the very first row
		columns := rows[0]

		// subslice rows in order to skip the first row (the columns)
		for _, row := range rows[1:len(rows)] {
			json, err := createRowToJSON(columns, row)
			if err != nil {
				log.Printf("ERROR: Could not join row with header - %s", err)
			} else {
				err := publish(rmqConn, os.Getenv("MESSAGEQUEUESERVER_EXCHANGE"), os.Getenv("MESSAGEQUEUESERVER_QUEUE"), string(json), true)
				if err != nil {
					log.Printf("ERROR: Could not publish message '%s' - %s", string(json), err)
				}
			}
		}
	}
}
