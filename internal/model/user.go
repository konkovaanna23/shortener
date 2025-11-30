package model

import (
	"fmt"
	"slices"
	"sync"
)

type UserURLS struct {
	mx       sync.RWMutex
	userURLS map[string][]string
}

func NewUserURLS() *UserURLS {
	return &UserURLS{
		userURLS: make(map[string][]string),
	}
}

func (lu *UserURLS) AddURLForUser(userID string, url string) {
	if userID != "" {
		lu.mx.Lock()
		defer lu.mx.Unlock()
		if !slices.Contains(lu.userURLS[userID], url) {
			lu.userURLS[userID] = append(lu.userURLS[userID], url)
		}
	}
}

func (lu *UserURLS) InitUsers(userurls map[string][]string) {
	lu.mx.Lock()
	defer lu.mx.Unlock()
	fmt.Println("инициализация user")
	fmt.Println(userurls)
	for key, value := range userurls {
		lu.userURLS[key] = value
	}
}

func (lu *UserURLS) GetURLsForUser(userID string) []string {
	if userID != "" {
		lu.mx.RLock()
		defer lu.mx.RUnlock()
		return lu.userURLS[userID]
	} else {
		return nil
	}
}

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
