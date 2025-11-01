package model_test

import (
	"testing"

	"github.com/konkovaanna23/shortener/internal/model"
)

func TestStorage_Get(t *testing.T) {
	s := model.NewStorage(5)
	sourceURL := "http://ya.ru"
	result := s.Add(sourceURL)
	resultURL := s.Get(result)
	if sourceURL != resultURL {
		t.Errorf("Get() = %v, want %v", resultURL, sourceURL)
	}
}

func TestStorage_Add(t *testing.T) {
	s := model.NewStorage(5)
	sourceURL := "http://ya.ru"
	result := s.Add(sourceURL)
	if len(result) != 5 {
		t.Errorf("Add() = %v, want %v", result, "random string 5 symbols")
	}
}
