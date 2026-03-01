package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/konkovaanna23/shortener/internal/config"
	"github.com/konkovaanna23/shortener/internal/config/db"
	"github.com/konkovaanna23/shortener/internal/grpc_server"
	"github.com/konkovaanna23/shortener/internal/handler"
	"github.com/konkovaanna23/shortener/internal/service"
	ss "github.com/konkovaanna23/shortener/pkg/shortenerservice"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {

	config.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGQUIT)
	defer signal.Stop(quitCh)

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
	server := handler.NewServer(cfg.URLserver, converter, cfg.Key, cfg.AuditFilePath, cfg.AuditURL, cfg.EnableHTTPS, cfg.TrustedSubnet)
	grpcServer := grpc_server.NewGrpcServer(converter, cfg.AuditFilePath, cfg.AuditURL)

	go func() {
		logrus.Printf("Сервер запущен на: %s", cfg.URLserver)
		if err := server.Start(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	if cfg.GrpcServer != "" {
		go func() {
			if err := StartGrpcServer(cfg.GrpcServer, grpcServer); err != nil {
				logrus.Error("Ошибка запуска gRPC сервера:", err)
			}
		}()
	}

	for {
		select {
		case <-quitCh:
			fmt.Fprintln(os.Stderr, "Получен сигнал SIGQUIT: goroutine dump (pprof)")

			if p := pprof.Lookup("goroutine"); p != nil {
				err := p.WriteTo(os.Stderr, 2)
				if err != nil {
					logrus.Error(err)
				}
			} else {
				fmt.Fprintln(os.Stderr, "pprof.Lookup(\"goroutine\") вернул nil")
			}
			continue

		case <-ctx.Done():
			logrus.Println("Сервер остановлен")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err := pprofSrv.Shutdown(shutdownCtx)
			if err != nil {
				logrus.Error(err)
			}
			return
		}
	}
}

func StartGrpcServer(host string, srv *grpc_server.GrpcServer) error {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_server.UnaryInterceptor))

	ss.RegisterShortenerServiceServer(s, srv)

	logrus.Info("сервер gRPC начал работу, host:", host)

	if err := s.Serve(listen); err != nil {
		return err
	}
	return nil
}
