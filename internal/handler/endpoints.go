package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/sirupsen/logrus"
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
	flagConflictError := false
	if err != nil {
		flagConflictError = errors.Is(err, model.ErrorConflictURL)
		if !flagConflictError {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
	logrus.Info("POST Сокращенный URL:", result)
	w.Header().Set("Content-Type", "text/plain")
	if flagConflictError {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

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
	flagConflictError := false
	if err != nil {
		flagConflictError = errors.Is(err, model.ErrorConflictURL)
		if !flagConflictError {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
	logrus.Info("POST Сокращенный URL:", result.URLShort)
	bodyResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if flagConflictError {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(bodyResult)

}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	err := s.converter.PingDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) newJSONBatchURL(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения BODY", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var urlDescription []*model.DescriptionURL
	err = json.Unmarshal(bodyBytes, &urlDescription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("POST Заданный JSON:", string(bodyBytes))
	result, err := s.converter.AddURLForBatch(urlDescription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bodyResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("POST Список сокращенных URL:", string(bodyResult))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(bodyResult)

}
