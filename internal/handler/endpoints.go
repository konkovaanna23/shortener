package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/konkovaanna23/shortener/internal/handler/audit"
	"github.com/konkovaanna23/shortener/internal/handler/middleware"
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
	user, _ := middleware.GetUserID(r.Context())
	logrus.Info("POST Заданный URL:", text, " заданный user:", user)
	result, err := s.converter.AddURL(text, user)
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

	s.auditor.Publish(audit.Event{Time: time.Now(),
		Action: "shorten",
		UserID: user,
		URL:    text})

}

func (s *Server) getURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shorturl")
	if shortURL == "" {
		http.Error(w, "Пустой URL", http.StatusBadRequest)
	}
	logrus.Info("GET Заданный URL:", shortURL)
	sourceURL, err := s.converter.GetURL(shortURL)
	if err != nil {
		if errors.Is(err, model.ErrorDeletedURL) {
			w.WriteHeader(http.StatusGone)
			return
		} else {
			logrus.Println("Ошибка:" + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	logrus.Info("GET Исходный URL:", sourceURL)
	s.auditor.Publish(audit.Event{Time: time.Now(),
		Action: "follow",
		URL:    sourceURL,
	})
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
	user, _ := middleware.GetUserID(r.Context())
	logrus.Info("POST JSONUrl Заданный URL:", urlRequest.URL, " заданный user:", user)
	result, err := s.converter.AddURLForRequest(urlRequest, user)
	flagConflictError := false
	if err != nil {
		flagConflictError = errors.Is(err, model.ErrorConflictURL)
		if !flagConflictError {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
	logrus.Info("POST JSONUrl Сокращенный URL:", result.URLShort)
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
	s.auditor.Publish(audit.Event{Time: time.Now(),
		Action: "shorten",
		UserID: user,
		URL:    urlRequest.URL})

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
	user, _ := middleware.GetUserID(r.Context())
	logrus.Info("POST Batch Заданный JSON:", string(bodyBytes), " заданный user:", user)
	result, err := s.converter.AddURLForBatch(urlDescription, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bodyResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("POST Batch Список сокращенных URL:", string(bodyResult))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(bodyResult)

}

func (s *Server) getURLForUser(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserID(r.Context())
	logrus.Info("GET UrlForUser Заданный user:", user)
	result, err := s.converter.GetURLsForUser(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	bodyResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("GET UrlForUser Список сокращенных URL:", string(bodyResult))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyResult)

}

func (s *Server) deleteURLForUser(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserID(r.Context())
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения BODY", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var urls []string
	err = json.Unmarshal(bodyBytes, &urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("Delete UrlForUser ", string(bodyBytes), " Заданный user:", user)
	go s.converter.DeleteURLsForUser(urls, user)
	w.WriteHeader(http.StatusAccepted)
}
