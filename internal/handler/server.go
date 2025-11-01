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

func NewServer(url, urlForShort string) *Server {

	mux := chi.NewRouter()

	s := &Server{
		mux:       mux,
		url:       url,
		converter: service.NewConverter(urlForShort),
	}
	s.mux.Post("/", s.newURL)
	s.mux.Get("/{shorturl}", s.getURL)
	return s
}

func (s *Server) Start() error {
	err := http.ListenAndServe(s.url, s.mux)
	return err
}
