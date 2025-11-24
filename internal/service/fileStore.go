package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/konkovaanna23/shortener/internal/file"
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/sirupsen/logrus"
)

func (c *Converter) GetInfoURLFromFile() (map[string]string, error) {
	var resultMap map[string]string
	data, err := file.ReadFromFile(c.filePath)
	if err != nil {
		logrus.Errorln("ошибка получения URLs из файла:", err)
		return nil, err
	} else {
		resultMap, err = c.decodeDataToOriginalMap(data)
		if err != nil {
			logrus.Errorln("ошибка перкодирования в map:", err)
			return nil, err
		}

	}
	return resultMap, nil
}

func (c *Converter) decodeDataToOriginalMap(data []byte) (map[string]string, error) {
	var result map[string]string
	if len(data) == 0 {
		return nil, errors.New("данные не переданы")
	} else {
		var urls []*model.DescriptionURL
		if err := json.Unmarshal(data, &urls); err != nil {
			return nil, err
		}
		result = c.convertDescriptionURLToOriginalMap(urls)
	}
	return result, nil
}

func (c *Converter) decodeDataToShortMap(data []byte) (map[string]string, error) {
	var result map[string]string
	if len(data) == 0 {
		return nil, errors.New("данные не переданы")
	} else {
		var urls []*model.DescriptionURL
		if err := json.Unmarshal(data, &urls); err != nil {
			return nil, err
		}
		result = c.convertDescriptionURLToShortMap(urls)
	}
	return result, nil
}

func (c *Converter) StoreURLInFile(shortURL, originalURL string) (string, error) {
	c.fMx.Lock()
	defer c.fMx.Unlock()
	resultMap, err := c.GetInfoURLFromFile()
	if err != nil {
		return "", err
	}
	shortURLCurrent, ok := resultMap[originalURL]
	if ok {
		return shortURLCurrent, model.ErrorConflictURL
	} else {
		resultMap[originalURL] = shortURL
	}
	data, err := c.encodeMapToData(resultMap)
	if err != nil {
		return "", err
	}
	err = file.SaveToFile(c.filePath, data)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (c *Converter) GetOriginalURLFromFile(shortURL string) (string, error) {
	c.fMx.RLock()
	defer c.fMx.RUnlock()
	var resultMap map[string]string
	data, err := file.ReadFromFile(c.filePath)
	if err != nil {
		logrus.Errorln("ошибка получения URLs из файла:", err)
		return "", err
	} else {
		resultMap, err = c.decodeDataToShortMap(data)
		if err != nil {
			logrus.Errorln("ошибка перкодирования в map:", err)
			return "", err
		}
		if originalURL, ok := resultMap[shortURL]; ok {
			return originalURL, nil
		} else {
			return "", fmt.Errorf("не существует оригинально URL для %s", shortURL)
		}

	}
}

// Возвращает мапу с ключами - оригинальные URL, значениями - короткие URL
func (c *Converter) convertDescriptionURLToOriginalMap(urls []*model.DescriptionURL) map[string]string {
	result := make(map[string]string)
	for _, desc := range urls {
		result[desc.Original] = desc.Short
	}
	return result
}

// Возвращает мапу с ключами - короткие URL, значениями - оригинальные URL
func (c *Converter) convertDescriptionURLToShortMap(urls []*model.DescriptionURL) map[string]string {
	result := make(map[string]string)
	for _, desc := range urls {
		result[desc.Short] = desc.Original
	}
	return result
}

func (c *Converter) encodeMapToData(mapURL map[string]string) ([]byte, error) {
	if len(mapURL) != 0 {
		list := make([]*model.DescriptionURL, len(mapURL))
		i := 0
		for key, value := range mapURL {
			list[i] = &model.DescriptionURL{Short: value, Original: key, ID: i + 1}
			i++
		}
		return json.Marshal(list)
	}
	return nil, nil
}

func (c *Converter) StoreURLsInFile(urls []*model.DescriptionURL) error {
	c.fMx.Lock()
	defer c.fMx.Unlock()
	resultMap, err := c.GetInfoURLFromFile()
	for _, url := range urls {
		if value, ok := resultMap[url.Original]; ok {
			url.Short = value
		} else {
			resultMap[url.Original] = url.Short
		}
	}
	data, err := c.encodeMapToData(resultMap)
	if err != nil {
		return err
	}
	err = file.SaveToFile(c.filePath, data)
	if err != nil {
		return err
	}
	return nil
}
