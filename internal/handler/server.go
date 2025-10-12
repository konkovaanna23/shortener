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

func NewServer(url string) *Server {

	mux := chi.NewRouter()

	s := &Server{
		mux:       mux,
		url:       url,
		converter: service.NewConverter("http://" + url),
	}
	s.mux.Post("/", s.newURL)
	s.mux.Get("/{shorturl}", s.getURL)
	return s
}

func (s *Server) Start() {
	if err := http.ListenAndServe(s.url, s.mux); err != nil {
		panic(err)
	}
}
