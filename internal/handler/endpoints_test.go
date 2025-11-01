package handler

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_newURL(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:8080/", strings.NewReader("https://yandex.ru"))
	req.Header.Set("Content-Type", "text/plain")

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	s := NewServer("localhost:8080", "http://localhost:8080")
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
	sourceURL := "https://google.com"
	s := NewServer("localhost:8080", "http://localhost:8080")
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
