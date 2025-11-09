package main

import (
	"context"
	"github.com/konkovaanna23/shortener/internal/config"
	"github.com/konkovaanna23/shortener/internal/file"
	"github.com/konkovaanna23/shortener/internal/handler"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetConfig()
	converter := service.NewConverter(cfg.URLforShort, cfg.FilePath)
	server := handler.NewServer(cfg.URLserver, converter)

	go func() {
		logrus.Printf("Сервер запущен на: %s", cfg.URLserver)
		if err := server.Start(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	<-sigChan

	data, err := converter.GetAllData()
	if err != nil {
		logrus.Error("Ошибка при получении данных для сохранения:", err)
	} else {
		if err := file.SaveToFile(cfg.FilePath, data); err != nil {
			logrus.Error("Ошибка сохранения в файл:", err)
		} else {
			logrus.Println("Данные сохранены в файл:", cfg.FilePath)
		}
	}
	logrus.Println("Сервер остановлен")

}
