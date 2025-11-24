package service

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/konkovaanna23/shortener/internal/config/db"
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/sirupsen/logrus"
)

const (
	lengthURL = 6
)

const (
	MODE_STORE_DB      = "DB"
	MODE_STORE_FILE    = "FILE"
	MODE_STORE_STORAGE = "STORAGE"
)

type Converter struct {
	url       string
	storage   *model.Storage
	db        *sqlx.DB
	modeStore string
	filePath  string
	fMx       sync.RWMutex
}

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	URLShort string `json:"result"`
}

func NewConverter(serverURL string, filePath string, db *sqlx.DB) *Converter {
	cvrt := &Converter{
		url:      serverURL,
		storage:  model.NewStorage(lengthURL),
		db:       db,
		filePath: filePath,
	}

	if db != nil {
		cvrt.modeStore = MODE_STORE_DB
		logrus.Println("Установлен режим сохранения в БД")
	} else {
		if filePath != "" {
			cvrt.modeStore = MODE_STORE_FILE
			logrus.Println("Установлен режим сохранения в файл")
		} else {
			cvrt.modeStore = MODE_STORE_STORAGE
			logrus.Println("Установлен режим сохранения в хранилище")
		}
	}
	return cvrt
}

func (c *Converter) AddURL(url string) (string, error) {
	if !c.isValidURL(url) {
		msg := fmt.Sprintf("URL [%s] не является валидным", url)
		return "", fmt.Errorf("%s", msg)
	}
	shortURL := c.storage.GenerateShortURL()
	var err error
	switch c.modeStore {
	case MODE_STORE_DB:
		if shortURL, err = c.StoreURLInDB(shortURL, url); err != nil {
			logrus.Errorln("ошибка сохранения в базу:", err)
		}
	case MODE_STORE_FILE:
		if shortURL, err = c.StoreURLInFile(shortURL, url); err != nil {
			logrus.Errorln("ошибка сохранения в файл:", err)
		}
	case MODE_STORE_STORAGE:
		shortURL, err = c.storage.Add(url, shortURL)
	}
	return c.url + "/" + shortURL, err
}

func (c *Converter) GetURL(shortURL string) (string, error) {
	var URL string
	ok := true
	var err error
	switch c.modeStore {
	case MODE_STORE_DB:
		URL, err = c.GetOriginalURLFromDB(shortURL)
		if err != nil {
			logrus.Errorln(err)
			ok = false
		}
	case MODE_STORE_FILE:
		URL, err = c.GetOriginalURLFromFile(shortURL)
		if err != nil {
			logrus.Errorln(err)
			ok = false
		}
	case MODE_STORE_STORAGE:
		URL, ok = c.storage.Get(shortURL)
	}
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

func (c *Converter) AddURLForRequest(url *URLRequest) (*URLResponse, error) {
	if url == nil {
		return nil, errors.New("передана пустая структура")
	}
	result, err := c.AddURL(url.URL)
	if err != nil {
		if errors.Is(err, model.ErrorConflictURL) {
			return &URLResponse{URLShort: result}, err
		} else {
			return nil, err
		}
	}
	return &URLResponse{URLShort: result}, nil
}

func (c *Converter) PingDB() error {
	err := db.Ping(c.db)
	if err != nil {
		return err
	}
	return nil
}

func (c *Converter) AddURLForBatch(urls []*model.DescriptionURL) ([]*model.DescriptionURL, error) {
	var errs []error
	for _, url := range urls {
		if !c.isValidURL(url.Original) {
			errs = append(errs, fmt.Errorf("URL [%s] не является валидным", url.Original))
		}
	}
	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	for _, url := range urls {
		url.Short = c.storage.GenerateShortURL()
	}

	switch c.modeStore {
	case MODE_STORE_DB:
		err := c.TranStoreURLInDB(urls)
		if err != nil {
			logrus.Error("ошибка сохранения в базу:", err)
		}
	case MODE_STORE_FILE:
		err := c.StoreURLsInFile(urls)
		if err != nil {
			logrus.Error("ошибка сохранения в базу:", err)
		}
	case MODE_STORE_STORAGE:
		for _, url := range urls {
			short, _ := c.storage.Add(url.Original, url.Short)
			url.Short = short
		}
	}
	for _, url := range urls {
		url.Original = ""
		url.Short = c.url + "/" + url.Short
	}
	return urls, nil
}
