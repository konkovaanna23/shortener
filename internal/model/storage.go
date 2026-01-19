package model

import (
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"slices"
	"sync"
)

var (
	ErrorConflictURL = errors.New("conflict url")
	ErrorDeletedURL  = errors.New("deleted url")
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lengthLetters = byte(len(letters))

// Storage - хранилище URL
type Storage struct {
	length     int
	urls       sync.Map
	deleteURLs sync.Map
	bufPool    sync.Pool
}

// NewStorage - конструктор хранилища.
func NewStorage(length int) *Storage {
	return &Storage{
		length:     length,
		urls:       sync.Map{}, //map[shortURL]originalURL
		deleteURLs: sync.Map{},
		bufPool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 8)
			},
		},
	}
}

// InitStorage - инициализация хранилища
// Параметры:
// - urls: map[shortURL]originalURL
func (s *Storage) InitStorage(urls map[string]string) {
	for key, value := range urls {
		s.urls.Store(key, value)
	}
}

// Add - добавление URL в хранилище.
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

// Delete - удаление URL из хранилища
func (s *Storage) Delete(short string) {
	_, ok := s.urls.Load(short)
	if ok {
		s.deleteURLs.Store(short, true)
	}
}

func (s *Storage) randomString(letters string) string {
	buf := s.bufPool.Get().([]byte)
	if len(buf) < s.length {
		buf = make([]byte, s.length)
	}
	defer s.bufPool.Put(buf)

	_, _ = cryptorand.Read(buf[:s.length])
	for i := 0; i < s.length; i++ {
		buf[i] = letters[buf[i]%lengthLetters]
	}
	return string(buf[:s.length])
}

// GenerateShortURL - генерация короткого URL
func (s *Storage) GenerateShortURL() string {
	return s.randomString(letters)
}

// Get - получение URL по короткому URL.
func (s *Storage) Get(shortURL string) (string, error) {
	URL, ok := s.urls.Load(shortURL)
	if !ok {
		msg := fmt.Sprintf("URL по короткому URL [%s] не существует", shortURL)
		return "", fmt.Errorf("%s", msg)
	}
	_, ok = s.deleteURLs.Load(shortURL)
	if ok {
		return "", ErrorDeletedURL
	}
	return URL.(string), nil
}

// GetAllURLMap - получение всех URL из хранилища.
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

// GetURLMapForList - получение URL из хранилища по списку коротких URL.
func (s *Storage) GetURLMapForList(list []string) map[string]string {
	resultMap := make(map[string]string)
	s.urls.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if slices.Contains(list, k) {
			_, ok := s.deleteURLs.Load(k)
			if !ok {
				resultMap[k] = v
			}
		}
		return true
	})
	return resultMap
}
