package main

import (
	"log"
	"os"
)

func main() {
	storage := NewPostgresStore()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("no port env set")
	}
	api := NewApi(storage, ":"+port)
	api.Run()
}
