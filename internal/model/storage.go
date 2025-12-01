package model

import (
	"errors"
	"math/rand"
	"slices"
	"sync"
)

var (
	ErrorConflictURL = errors.New("conflict url")
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Storage struct {
	length int
	urls   sync.Map
}

func NewStorage(length int) *Storage {
	return &Storage{
		length: length,
		urls:   sync.Map{}, //map[shortURL]originalURL
	}
}

func (s *Storage) InitStorage(urls map[string]string) {
	for key, value := range urls {
		s.urls.Store(key, value)
	}
}

func (s *Storage) Add(url string, short string) (string, error) {
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
		s.urls.Store(short, url)
		return short, nil
	} else {
		return shortURL, ErrorConflictURL
	}

}

func (s *Storage) randomString(letters string) string {
	b := make([]byte, s.length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *Storage) GenerateShortURL() string {
	return s.randomString(letters)
}

func (s *Storage) Get(shortURL string) (string, bool) {
	URL, ok := s.urls.Load(shortURL)
	if !ok {
		return "", ok
	}
	return URL.(string), ok
}

func (s *Storage) GetAllURLMap() map[string]string {
	resultMap := make(map[string]string)
	s.urls.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		resultMap[k] = v
		return true
	})
	return resultMap
}

func (s *Storage) GetURLMapForList(list []string) map[string]string {
	resultMap := make(map[string]string)
	s.urls.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if slices.Contains(list, k) {
			resultMap[k] = v
		}
		return true
	})
	return resultMap
}
