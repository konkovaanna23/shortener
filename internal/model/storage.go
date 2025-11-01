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
	shortURL := ""
	s.urls.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if v == url {
			shortURL = k
			return false
		}
		return true
	})
	if shortURL == "" {
		str := s.RandomString(letters)
		s.urls.Store(str, url)
		return str
	} else {
		return shortURL
	}

}

func (s *Storage) RandomString(letters string) string {
	b := make([]byte, s.length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *Storage) Get(shortURL string) (string, bool) {
	URL, ok := s.urls.Load(shortURL)
	if !ok {
		return "", ok
	}
	return URL.(string), ok
}
