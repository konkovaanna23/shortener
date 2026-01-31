package model

import (
	"slices"
	"sync"
)

// UserURLS хранит URL-ы пользователя.
type UserURLS struct {
	mx       sync.RWMutex
	userURLS map[string][]string
}

// NewUserURLS создает новый экземпляр UserURLS.
func NewUserURLS() *UserURLS {
	return &UserURLS{
		userURLS: make(map[string][]string),
	}
}

// AddURLForUser добавляет URL для пользователя.
func (lu *UserURLS) AddURLForUser(userID string, url string) {
	if userID != "" {
		lu.mx.Lock()
		defer lu.mx.Unlock()
		if !slices.Contains(lu.userURLS[userID], url) {
			lu.userURLS[userID] = append(lu.userURLS[userID], url)
		}
	}
}

// InitUsers инициализирует пользователей.
func (lu *UserURLS) InitUsers(userurls map[string][]string) {
	lu.mx.Lock()
	defer lu.mx.Unlock()
	for key, value := range userurls {
		lu.userURLS[key] = value
	}
}

// GetUserURLS возвращает URL-ы пользователя.
func (lu *UserURLS) GetURLsForUser(userID string) []string {
	if userID != "" {
		lu.mx.RLock()
		defer lu.mx.RUnlock()
		return lu.userURLS[userID]
	} else {
		return nil
	}
}

// ExistURLForUser проверяет, существует ли URL для пользователя.
func (lu *UserURLS) ExistURLForUser(userID string, url string) bool {
	if userID == "" {
		return false
	}
	lu.mx.RLock()
	defer lu.mx.RUnlock()
	for _, u := range lu.userURLS[userID] {
		if u == url {
			return true
		}
	}
	return false
}
