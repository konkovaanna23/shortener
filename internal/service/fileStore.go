package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/konkovaanna23/shortener/internal/file"
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/sirupsen/logrus"
)

func (c *Converter) getInfoURLFromFile() ([]*model.DescriptionURL, error) {
	var urls []*model.DescriptionURL
	data, err := file.ReadFromFile(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("файл %s не найден", c.filePath)
			return nil, nil
		} else {
			logrus.Errorln("ошибка получения URLs из файла:", err)
			return nil, err
		}
	} else {

		if err := json.Unmarshal(data, &urls); err != nil {
			return nil, err
		}
	}
	return urls, nil
}

// StoreURLInFile функция сохранения URL в файл.
func (c *Converter) StoreURLInFile(shortURL, originalURL, user string) (string, error) {
	c.fMx.Lock()
	defer c.fMx.Unlock()
	sourceURLs, err := c.getInfoURLFromFile()
	if err != nil {
		return "", err
	}
	if sourceURLs == nil {
		sourceURLs = make([]*model.DescriptionURL, 0)
	}
	for _, url := range sourceURLs {
		if url.Original == originalURL {
			return url.Short, model.ErrorConflictURL
		}
	}
	sourceURLs = append(sourceURLs, &model.DescriptionURL{Short: shortURL, Original: originalURL, ID: len(sourceURLs) + 1, UserID: user})
	data, err := json.Marshal(sourceURLs)
	if err != nil {
		return "", err
	}
	err = file.SaveToFile(c.filePath, data)
	if err != nil {
		return "", err
	}
	logrus.Infof("сохранен URL=%s для пользователя %s в файл %s", shortURL, user, c.filePath)
	return shortURL, nil
}

// GetOriginalURLFromFile функция получения оригинального URL по короткому из файла.
func (c *Converter) GetOriginalURLFromFile(shortURL string) (string, error) {
	c.fMx.RLock()
	defer c.fMx.RUnlock()
	sourceURLs, err := c.getInfoURLFromFile()
	if err != nil {
		return "", err
	}
	if sourceURLs == nil {
		return "", nil
	}
	for _, url := range sourceURLs {
		if url.Short == shortURL {
			if url.IsDeleted != nil && *url.IsDeleted {
				return "", model.ErrorDeletedURL
			}
			return url.Original, nil
		}
	}
	return "", fmt.Errorf("не существует оригинально URL для %s", shortURL)
}

// StoreURLsInFile функция сохранения нескольких URL в файл.
func (c *Converter) StoreURLsInFile(urls []*model.DescriptionURL, user string) error {
	c.fMx.Lock()
	defer c.fMx.Unlock()
	sourceURLs, err := c.getInfoURLFromFile()
	if err != nil {
		return err
	}
	if sourceURLs == nil {
		sourceURLs = make([]*model.DescriptionURL, 0)
	}
	for _, url := range urls {
		flagNew := true
		for _, sourceURL := range sourceURLs {
			if sourceURL.Original == url.Original {
				url.Short = sourceURL.Short
				flagNew = false
				break
			}
		}
		if flagNew {
			url.UserID = user
			sourceURLs = append(sourceURLs, url)
		}
	}
	data, err := json.Marshal(sourceURLs)
	if err != nil {
		return err
	}
	err = file.SaveToFile(c.filePath, data)
	if err != nil {
		return err
	}
	return nil
}

// GetInfoUserURLFromFile функция получения информации о URL пользователя из файла.
func (c *Converter) GetInfoUserURLFromFile(user string) ([]*model.DescriptionURL, error) {
	c.fMx.RLock()
	defer c.fMx.RUnlock()
	sourceURLs, err := c.getInfoURLFromFile()
	if err != nil {
		return nil, err
	}
	if sourceURLs == nil {
		return nil, nil
	}
	urlForUser := make([]*model.DescriptionURL, 0)
	for _, url := range sourceURLs {
		if url.UserID == user && (url.IsDeleted == nil || (url.IsDeleted != nil && !*url.IsDeleted)) {
			u := url.Copy()
			u.UserID = ""
			u.ID = 0
			u.Correlation = ""
			urlForUser = append(urlForUser, u)
		}
	}
	return urlForUser, nil
}

// DeleteURLsFromFile функция удаления URL из файла.
func (c *Converter) DeleteURLsFromFile(urls []*model.DescriptionURL) error {
	c.fMx.Lock()
	defer c.fMx.Unlock()
	sourceURLs, err := c.getInfoURLFromFile()
	if err != nil {
		return err
	}
	if sourceURLs == nil {
		sourceURLs = make([]*model.DescriptionURL, 0)
	}
	for _, url := range urls {
		flagFind := false
		for _, sourceURL := range sourceURLs {
			if sourceURL.Short == url.Short && sourceURL.UserID == url.UserID {
				flagTrue := true
				sourceURL.IsDeleted = &flagTrue
				flagFind = true
				break
			}
		}
		if !flagFind {
			logrus.Infof("Не существует URL=%s у пользователя %s", url.Short, url.UserID)
		}
	}
	data, err := json.Marshal(sourceURLs)
	if err != nil {
		return err
	}
	err = file.SaveToFile(c.filePath, data)
	if err != nil {
		return err
	}
	return nil
}
