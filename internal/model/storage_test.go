package model_test

import (
	"testing"

	"github.com/konkovaanna23/shortener/internal/model"
)

func TestStorage_Get(t *testing.T) {
	s := model.NewStorage(5)
	sourceURL := "http://ya.ru"
	shortURL := s.GenerateShortURL()
	result, _ := s.Add(sourceURL, shortURL)
	resultURL, ok := s.Get(result)
	if ok != true {
		t.Errorf("URL не найден")
	}
	if sourceURL != resultURL {
		t.Errorf("Get() = %v, want %v", resultURL, sourceURL)
	}
}

func TestStorage_Add(t *testing.T) {
	s := model.NewStorage(5)
	sourceURL := "http://ya.ru"
	shortURL := s.GenerateShortURL()
	result, _ := s.Add(sourceURL, shortURL)
	if result != shortURL {
		t.Errorf("Get() = %v, want %v", result, shortURL)
	}
}
