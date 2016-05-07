package main

import (
	"bufio"
	"io"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"os"
)

// AwsInfo stores all aws configuration.
type AwsInfo struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
}

// Uses AwsInfo to connect to AWS, need to fix hard-coded region
func awsConnect(i *AwsInfo) *s3.Bucket {
	AWSAuth := aws.Auth{
		AccessKey: i.AccessKey,
		SecretKey: i.SecretKey,
	}
	region := aws.USEast
	connection := s3.New(AWSAuth, region)

	return connection.Bucket(i.Bucket)
}

func awsDownload(bucket *s3.Bucket, remoteFile string, localFile string) {

	downloadBytes, err := bucket.Get(remoteFile)
	if err != nil {
		log.Printf("ERROR: could not download %s from S3 - %s", remoteFile, err)
	} else {

		downloadFile, err := os.Create(localFile)
		if err != nil {
			log.Printf("ERROR: could not create %s - %s", localFile, err)
		} else {

			defer downloadFile.Close()

			downloadBuffer := bufio.NewWriter(downloadFile)
			defer downloadBuffer.Flush()
			downloadBuffer.Write(downloadBytes)

			io.Copy(downloadBuffer, downloadFile)
			log.Printf("%s downloaded to %s", remoteFile, localFile)
		}
	}
}
