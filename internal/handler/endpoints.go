package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (s *Server) newURL(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения BODY", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	text := string(bodyBytes)
	logrus.Info("POST Заданный URL:", text)
	result, err := s.converter.AddURL(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("POST Сокращенный URL:", result)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))

}

func (s *Server) getURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shorturl")
	if shortURL == "" {
		http.Error(w, "Пустой URL", http.StatusBadRequest)
	}
	logrus.Info("GET Заданный URL:", shortURL)
	sourceURL, err := s.converter.GetURL(shortURL)
	if err != nil {
		logrus.Println("Ошибка:" + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("GET Исходный URL:", sourceURL)
	http.Redirect(w, r, sourceURL, http.StatusTemporaryRedirect)
}

func (s *Server) newJSONURL(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения BODY", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var urlRequest *service.URLRequest
	err = json.Unmarshal(bodyBytes, &urlRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("POST Заданный URL:", urlRequest.URL)
	result, err := s.converter.AddURLForRequest(urlRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("POST Сокращенный URL:", result.URLShort)
	bodyResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(bodyResult)

}
