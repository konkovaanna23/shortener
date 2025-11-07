package model

import (
	"sync"
)

type DescriptionURL struct {
	ID       int    `json:"uuid"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type ListURL struct {
	URLs []*DescriptionURL
	mx   sync.RWMutex
}

func NewListURL() *ListURL {
	return &ListURL{
		URLs: make([]*DescriptionURL, 0),
	}
}

func (lu *ListURL) AddItеm(url *DescriptionURL) {
	lu.mx.Lock()
	defer lu.mx.Unlock()
	if url != nil {
		idx := len(lu.URLs)
		url.ID = idx + 1
		lu.URLs = append(lu.URLs, url)
	}
}
