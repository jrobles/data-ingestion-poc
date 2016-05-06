package main

import (
	"time"
)

func main() {
	for {
		time.Sleep(100 * time.Second)
		print("hello")
	}
}
