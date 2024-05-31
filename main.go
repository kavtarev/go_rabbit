package main

import (
	"log"
	"os"
)

func main() {
	storage := NewPostgresStore()
	channel := NewChannel()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("no port env set")
	}
	api := NewApi(storage, ":"+port, channel)
	go Receiver()
	api.Run()
}
