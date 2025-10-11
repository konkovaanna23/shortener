package service

import (
	"github.com/konkovaanna23/shortener/internal/model"
)

const (
	LENGTH_URL = 6
)

type Converter struct {
	url     string
	storage *model.Storage
}

func NewConverter(serverUrl string) *Converter {
	return &Converter{
		url:     serverUrl,
		storage: model.NewStorage(LENGTH_URL),
	}
}

func (c *Converter) AddUrl(url string) string {
	result := c.storage.Add(url)
	return c.url + "/" + result
}

func (c *Converter) GetUrl(shortUrl string) string {
	return c.storage.Get(shortUrl)
}
