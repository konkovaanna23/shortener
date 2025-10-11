package main

import (
	"log"

	"github.com/konkovaanna23/shortener/internal/handler"
)

func main() {
	server := handler.NewServer("localhost:8080")
	log.Println("Север запущен на http://localhost:8080")
	server.Start()
}
