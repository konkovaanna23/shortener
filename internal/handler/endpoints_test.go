package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_newURL(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:8080/", strings.NewReader("https://google.com"))
	req.Header.Set("Content-Type", "text/plain")

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	s := NewServer("localhost:8080")
	s.newOrGetURL(recorder, req)
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
	req, err := http.NewRequest("POST", "http://localhost:8080/", strings.NewReader(sourceURL))
	req.Header.Set("Content-Type", "text/plain")

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	s := NewServer("localhost:8080")
	s.newOrGetURL(recorder, req)
	shortURL := recorder.Body.String()
	parts := strings.Split(shortURL, "/")
	var suffix string
	if len(parts) == 4 {
		suffix = parts[3]
	}
	req, err = http.NewRequest("GET", "http://localhost:8080/"+suffix, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder = httptest.NewRecorder()
	s.newOrGetURL(recorder, req)
	if recorder.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"expected status code %d, got %d",
			http.StatusTemporaryRedirect,
			recorder.Code,
		)
	}
	if recorder.Header().Get("Location") != sourceURL {
		t.Errorf(
			"expected Location header to be %s, got %s", sourceURL, recorder.Header().Get("Location"),
		)
	}
}
