// Package handler - модуль для работы с HTTP-запросами.
package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/konkovaanna23/shortener/internal/handler/audit"
	"github.com/konkovaanna23/shortener/internal/handler/middleware"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

// Server - структура сервера.
type Server struct {
	url           string
	mux           *chi.Mux
	converter     *service.Converter
	srv           *http.Server
	auditor       *audit.Publisher
	enableHTTPS   bool
	trustedSubnet string
}

// NewServer создаёт и настраивает новый экземпляр HTTP-сервера для сервиса сокращения URL.
//
// Инициализирует маршрутизацию, подключает middleware и настраивает аудит событий.
// Поддерживает:
//   - Сжатие запросов/ответов (gzip)
//   - Логирование всех запросов
//   - Аутентификацию с использованием HMAC-подписи (если ключ задан)
//   - Аудит POST-запросов через файл и/или внешний HTTP-сервис
//
// Аудит включается, если указан один или оба параметра: auditFile или auditURL.
// События аудита публикуются при создании и удалении URL.
//
// Параметры:
//   - url: адрес, на котором будет запущен сервер (например, ":8080")
//   - converter: бизнес-логика для генерации и хранения сокращённых URL
//   - key: секретный ключ для проверки HMAC-подписи (может быть пустым)
//   - auditFile: путь к файлу для записи аудит-событий (если пуст — не используется)
//   - auditURL: URL внешнего сервиса для отправки аудит-событий (если пуст — не используется)
//
// Возвращает указатель на *Server, готовый к запуску.
//
// Пример использования:
//
//	converter := service.NewConverter(...)
//	server := handler.NewServer(
//	    ":8080",
//	    converter,
//	    "my-secret-key",
//	    "/var/log/audit.log",
//	    "https://audit-service/api/events",
//	)
//
//	go server.Start(ctx)
func NewServer(url string, converter *service.Converter, key string, auditFile string, auditURL string, enableHTTPS bool, trustedSubnet string) *Server {

	mux := chi.NewRouter()

	s := &Server{
		mux:           mux,
		url:           url,
		converter:     converter,
		enableHTTPS:   enableHTTPS,
		trustedSubnet: trustedSubnet,
		auditor:       audit.NewAuditor(auditFile, auditURL),
	}
	s.mux.Use(middleware.CompressMiddleware)
	s.mux.Use(middleware.LoggingMiddleware)
	s.mux.Use(middleware.AuthMiddleware(key))
	s.mux.Post("/", s.newURL)
	s.mux.Get("/{shorturl}", s.getURL)
	s.mux.Get("/ping", s.ping)
	s.mux.Post("/api/shorten", s.newJSONURL)
	s.mux.Post("/api/shorten/batch", s.newJSONBatchURL)
	s.mux.Get("/api/user/urls", s.getURLForUser)
	s.mux.Delete("/api/user/urls", s.deleteURLForUser)
	if trustedSubnet != "" {
		s.mux.Get("/api/internal/stats", s.stats)
	}
	s.srv = &http.Server{
		Addr:    url,
		Handler: mux,
	}
	if enableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("shortener.ru", "www.shortener.ru"),
		}
		s.srv.TLSConfig = manager.TLSConfig()
	}

	return s
}

// Start запускает HTTP-сервер в отдельной горутине.
func (s *Server) Start(ctx context.Context) error {

	go func() {
		<-ctx.Done()
		if err := s.srv.Shutdown(ctx); err != nil {
			logrus.Error("Ошибка при закрытии сервера:", err)
		}
	}()

	if s.enableHTTPS {
		err := s.srv.ListenAndServeTLS("", "")
		return err
	} else {
		err := s.srv.ListenAndServe()
		return err
	}

}

func (s *Server) GetHandler() http.Handler {
	return s.mux
}
