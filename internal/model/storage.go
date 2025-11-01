package model

import (
	"math/rand"
	"sync"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Storage struct {
	length int
	urls   sync.Map
}

func NewStorage(length int) *Storage {
	return &Storage{
		length: length,
		urls:   sync.Map{},
	}
}

func (s *Storage) Add(url string) string {
	value, ok := s.urls.Load(url)
	if !ok {
		str := s.RandomString(letters)
		s.urls.Store(url, str)
		return str
	} else {
		return value.(string)
	}

}

func (s *Storage) RandomString(letters string) string {
	b := make([]byte, s.length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *Storage) Get(shortURL string) string {
	sourceURL := ""
	s.urls.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if v == shortURL {
			sourceURL = k
		}
		return true
	})
	return sourceURL
}
