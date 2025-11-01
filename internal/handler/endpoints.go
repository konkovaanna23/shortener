package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) newURL(w http.ResponseWriter, r *http.Request) {
	/*if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Некорректный Content-Type", http.StatusBadRequest)
			return
	}*/
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения BODY", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	text := string(bodyBytes)
	/*if text == "" {
		http.Error(w, "Пустой URL", http.StatusBadRequest)
		return
	}*/
	log.Println("POST Заданный URL:", text)
	result := s.converter.AddURL(text)
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
	sourceURL := s.converter.GetURL(shortURL)
	log.Println("GET Исходный URL:", sourceURL)
	w.Header().Set("Location", sourceURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
