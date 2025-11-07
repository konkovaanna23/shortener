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
	"sync"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.GetConfig()
	converter := service.NewConverter(cfg.URLforShort, cfg.FilePath)
	server := handler.NewServer(cfg.URLserver, converter)

	wg.Add(1)
	go func() {
		defer wg.Done()

		go func() {
			logrus.Println("Сервер запущен на :", cfg.URLserver)
			err := server.Start()
			if err != nil {
				logrus.Errorln(err)
			}
		}()
		<-ctx.Done()
		logrus.Println("Отмена контекста...")
	}()

	<-sigChan

	cancel()

	file.SaveToFile(cfg.FilePath, converter.GetAllData())
	logrus.Println("Сервер остановлен")
	wg.Wait()
}
