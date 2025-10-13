package main

import (
	"log"

	"github.com/konkovaanna23/shortener/internal/config"
	"github.com/konkovaanna23/shortener/internal/handler"
)

func main() {
	cfg := config.GetConfig()
	server := handler.NewServer(cfg.URLserver, cfg.URLforShort)
	log.Println("Сервер запущен на :", cfg.URLserver)
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
