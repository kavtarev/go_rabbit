package main

import (
	"log"
	"os"
)

func main23() {
	storage := NewPostgresStore()
	channel := NewChannel()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("no port env set")
	}
	api := NewApi(storage, ":"+port, channel)
	// go Receiver(1)
	// go Receiver(2)
	go ReceiverWithExchange()
	go ReceiverWithExchange()
	api.Run()
}
