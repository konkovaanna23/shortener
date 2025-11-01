package service

import (
	"github.com/konkovaanna23/shortener/internal/model"
)

const (
	lengthURL = 6
)

type Converter struct {
	url     string
	storage *model.Storage
}

func NewConverter(serverURL string) *Converter {
	return &Converter{
		url:     serverURL,
		storage: model.NewStorage(lengthURL),
	}
}

func (c *Converter) AddURL(url string) string {
	result := c.storage.Add(url)
	return c.url + "/" + result
}

func (c *Converter) GetURL(shortURL string) string {
	return c.storage.Get(shortURL)
}
