package main

import (
	"github.com/konkovaanna23/shortener/internal/config"
	"github.com/konkovaanna23/shortener/internal/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.GetConfig()
	server := handler.NewServer(cfg.URLserver, cfg.URLforShort)
	logrus.Println("Сервер запущен на :", cfg.URLserver)
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
