package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func csvImport(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(500)
	} else {
		file := Payload{}
		json.Unmarshal([]byte(body), &file)

		f := &AwsInfo{
			AccessKey: file.Key,
			SecretKey: file.Secret,
			Region:    file.Region,
			Bucket:    file.Bucket,
		}

		localFile := "/tmp/" + file.Filename + "." + file.Extension
		c := awsConnect(f)
		awsDownload(c, string(file.Path+"/"+file.Filename+"."+file.Extension), localFile)
	}
}
