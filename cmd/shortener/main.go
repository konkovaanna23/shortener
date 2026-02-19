package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/konkovaanna23/shortener/internal/config"
	"github.com/konkovaanna23/shortener/internal/config/db"
	"github.com/konkovaanna23/shortener/internal/handler"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/sirupsen/logrus"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {

	config.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// pprof server
	pprofSrv := &http.Server{Addr: "localhost:6060"}
	go func() {
		logrus.Println("pprof доступен на http://localhost:6060/debug/pprof/")
		if err := pprofSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Error("pprof server error: ", err)
		}
	}()

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

	converter := service.NewConverter(ctx, cfg.URLforShort, cfg.FilePath, database, cfg.BufferSize, cfg.BatchSize, cfg.TimeFlushDel)
	server := handler.NewServer(cfg.URLserver, converter, cfg.Key, cfg.AuditFilePath, cfg.AuditURL, cfg.EnableHTTPS)

	go func() {
		logrus.Printf("Сервер запущен на: %s", cfg.URLserver)
		if err := server.Start(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	<-ctx.Done()
	logrus.Println("Сервер остановлен")

	// graceful shutdown pprof
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = pprofSrv.Shutdown(shutdownCtx)
}
