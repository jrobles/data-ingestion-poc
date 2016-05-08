package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

		for _, row := range rows {
			fmt.Println(row)
		}

	}
}
