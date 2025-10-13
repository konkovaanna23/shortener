package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) newURL(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения BODY", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	text := string(bodyBytes)
	log.Println("POST Заданный URL:", text)
	result, err := s.converter.AddURL(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Println("POST Сокращенный URL:", result)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))

}

func (s *Server) getURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shorturl")
	if shortURL == "" {
		http.Error(w, "Пустой URL", http.StatusBadRequest)
	}
	log.Println("GET Заданный URL:", shortURL)
	sourceURL, err := s.converter.GetURL(shortURL)
	if err != nil {
		log.Println("Ошибка:" + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Println("GET Исходный URL:", sourceURL)
	w.Header().Set("Location", sourceURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
