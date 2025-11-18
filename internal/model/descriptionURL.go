package model

import (
	"sync"
)

type DescriptionURL struct {
	ID          int    `json:"uuid,omitempty"`
	Short       string `json:"short_url,omitempty"`
	Original    string `json:"original_url,omitempty"`
	Correlation string `json:"correlation_id,omitempty"`
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
