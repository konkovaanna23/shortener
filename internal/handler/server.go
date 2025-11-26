package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/konkovaanna23/shortener/internal/handler/middleware"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/sirupsen/logrus"
)

type Server struct {
	url       string
	mux       *chi.Mux
	converter *service.Converter
	srv       *http.Server
}

func NewServer(url string, converter *service.Converter) *Server {

	mux := chi.NewRouter()

	s := &Server{
		mux:       mux,
		url:       url,
		converter: converter,
	}
	s.mux.Use(middleware.CompressMiddleware)
	s.mux.Use(middleware.LoggingMiddleware)
	s.mux.Post("/", s.newURL)
	s.mux.Get("/{shorturl}", s.getURL)
	s.mux.Get("/ping", s.ping)
	s.mux.Post("/api/shorten", s.newJSONURL)
	s.mux.Post("/api/shorten/batch", s.newJSONBatchURL)
	s.srv = &http.Server{
		Addr:    url,
		Handler: mux,
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {

	go func() {
		<-ctx.Done()
		if err := s.srv.Shutdown(ctx); err != nil {
			logrus.Error("Ошибка при закрытии сервера:", err)
		}
	}()

	err := s.srv.ListenAndServe()
	return err
}
