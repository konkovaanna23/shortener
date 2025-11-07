package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/konkovaanna23/shortener/internal/service"
)

type Server struct {
	url       string
	mux       *chi.Mux
	converter *service.Converter
}

func NewServer(url string, converter *service.Converter) *Server {

	mux := chi.NewRouter()

	s := &Server{
		mux:       mux,
		url:       url,
		converter: converter,
	}
	s.mux.Post("/", s.LoggingMiddleware(http.HandlerFunc(s.newURL)))
	s.mux.Get("/{shorturl}", s.LoggingMiddleware(http.HandlerFunc(s.getURL)))
	s.mux.Post("/api/shorten", s.LoggingMiddleware(http.HandlerFunc(s.newJSONURL)))
	return s
}

func (s *Server) Start() error {
	err := http.ListenAndServe(s.url, s.mux)
	return err
}
