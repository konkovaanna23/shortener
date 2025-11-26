package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/konkovaanna23/shortener/internal/config/db"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestServer_newURL(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:8080/", strings.NewReader("https://yandex.ru"))
	req.Header.Set("Content-Type", "text/plain")

	if err != nil {
		t.Fatal(err)
	}
	database, _ := db.NewConnect("")
	converter := service.NewConverter("http://localhost:8080", "test.json", database)
	recorder := httptest.NewRecorder()
	s := NewServer("localhost:8080", converter)
	s.newURL(recorder, req)
	if recorder.Code != http.StatusCreated {
		t.Errorf(
			"expected status code %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}
	if recorder.Body.String() == "" {
		t.Errorf(
			"expected response body to be non-empty",
		)
	}
	if recorder.Header().Get("Content-Type") != "text/plain" {
		t.Errorf(
			"expected Content-Type header to be text/plain",
		)
	}

}

func TestServer_getURL(t *testing.T) {
	database, _ := db.NewConnect("")
	sourceURL := "https://google.com"
	converter := service.NewConverter("http://localhost:8080", "test.json", database)
	s := NewServer("localhost:8080", converter)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()
	response, err := http.Post(ts.URL+"/", "text/plain", strings.NewReader(sourceURL))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var shortURL string
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	shortURL = string(bodyBytes)
	parts := strings.Split(shortURL, "/")
	var suffix string
	if len(parts) == 4 {
		suffix = parts[3]
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	redirectResp, err := client.Get(ts.URL + "/" + suffix)
	if err != nil {
		t.Fatal(err)
	}
	defer redirectResp.Body.Close()
	assert.Equal(t, sourceURL, redirectResp.Header.Get("Location"))
	assert.Equal(t, http.StatusTemporaryRedirect, redirectResp.StatusCode)

}

func TestServer_newJsonURL(t *testing.T) {
	database, _ := db.NewConnect("")
	converter := service.NewConverter("http://localhost:8080", "test.json", database)
	s := NewServer("localhost:8080", converter)

	input := &service.URLRequest{
		URL: "https://example.com",
	}
	bodyBytes, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.newJSONURL(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Ожидался статус 201, получили %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Ожидался Content-Type application/json, получили %s", contentType)
	}

	var result map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Ответ не является валидным JSON: %v", err)
	}

	_, exists := result["result"]
	if !exists {
		t.Error("В JSON отсутствует поле 'result'")
	}

}
