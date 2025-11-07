package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/konkovaanna23/shortener/internal/file"
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/sirupsen/logrus"
)

const (
	lengthURL = 6
)

type Converter struct {
	url     string
	storage *model.Storage
}

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	URLShort string `json:"result"`
}

func NewConverter(serverURL string, filePath string) *Converter {
	cvrt := &Converter{
		url:     serverURL,
		storage: model.NewStorage(lengthURL),
	}
	if filePath != "" {
		data, err := file.ReadFromFile(filePath)
		if err != nil {
			logrus.Errorln(err)

		} else {
			resultMap, err := cvrt.decodeDataToMap(data)
			if err != nil {
				logrus.Errorln(err)
			} else {
				cvrt.storage.InitStorage(resultMap)
			}
		}

	}
	return cvrt
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

func (c *Converter) AddURLForRequest(url *URLRequest) (*URLResponse, error) {
	if url == nil {
		return nil, errors.New("передана пустая структура")
	}
	result, err := c.AddURL(url.URL)
	if err != nil {
		return nil, err
	}
	return &URLResponse{URLShort: result}, nil
}

func (c *Converter) decodeDataToMap(data []byte) (map[string]string, error) {
	result := make(map[string]string)
	if len(data) == 0 {
		return nil, errors.New("данные не переданы")
	} else {
		var urls []*model.DescriptionURL
		if err := json.Unmarshal(data, &urls); err != nil {
			return nil, err
		}
		for _, desc := range urls {
			result[desc.Short] = desc.Original
		}
	}
	return result, nil
}

func (c *Converter) encodeMapToData() ([]byte, error) {
	result := c.storage.GetAllURLMap()
	if len(result) != 0 {
		list := model.NewListURL()
		for key, value := range result {
			list.AddItеm(&model.DescriptionURL{Short: key, Original: value})
		}
		if data, err := json.Marshal(list.URLs); err != nil {
			return nil, err
		} else {
			return data, nil
		}
	}
	return nil, nil
}

func (c *Converter) GetAllData() []byte {
	if data, err := c.encodeMapToData(); err != nil {
		logrus.Errorln(err)
		return nil
	} else {
		return data
	}
}
