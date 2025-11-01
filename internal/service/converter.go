package service

import (
	"fmt"
	"github.com/konkovaanna23/shortener/internal/model"
	"net/url"
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

func (c *Converter) AddURL(url string) (string, error) {
	if !c.isValidURL(url) {
		msg := fmt.Sprintf("URL [%s] не является валидным", url)
		return "", fmt.Errorf("%s", msg)
	}
	result := c.storage.Add(url)
	return c.url + "/" + result, nil
}

func (c *Converter) GetURL(shortURL string) (string, error) {
	URL, ok := c.storage.Get(shortURL)
	if !ok {
		msg := fmt.Sprintf("URL по короткому URL [%s] не существует", shortURL)
		return "", fmt.Errorf("%s", msg)
	}
	return URL, nil
}

func (c *Converter) isValidURL(input string) bool {
	u, err := url.Parse(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
