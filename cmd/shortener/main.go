package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/konkovaanna23/shortener/internal/config"
	"github.com/konkovaanna23/shortener/internal/config/db"
	"github.com/konkovaanna23/shortener/internal/handler"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.GetConfig()
	database, err := db.NewConnect(cfg.DSN)
	if err != nil {
		logrus.Error("Ошибка при подключении к базе данных:", err)
	} else {
		logrus.Println("Подключение к базе данных успешно")
		if err := db.RunMigrations(cfg.DSN); err != nil {
			logrus.Error("Ошибка при установке миграций:", err)
			database = nil
		}
	}

	converter := service.NewConverter(ctx, cfg.URLforShort, cfg.FilePath, database)
	server := handler.NewServer(cfg.URLserver, converter)

	go func() {
		logrus.Printf("Сервер запущен на: %s", cfg.URLserver)
		if err := server.Start(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	<-ctx.Done()
	logrus.Println("Сервер остановлен")

}
