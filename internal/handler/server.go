package handler

import (
	"github.com/konkovaanna23/shortener/internal/service"
	"net/http"
)

type Server struct {
	url       string
	mux       *http.ServeMux
	converter *service.Converter
}

func NewServer(url string) *Server {

	mux := http.NewServeMux()
	s := &Server{
		mux:       mux,
		url:       url,
		converter: service.NewConverter("http://" + url),
	}
	mux.HandleFunc("/", s.newOrGetUrl)
	return s
}

func (s *Server) Start() {
	if err := http.ListenAndServe(s.url, s.mux); err != nil {
		panic(err)
	}
}
