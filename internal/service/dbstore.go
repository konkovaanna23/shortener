package service

import (
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/sirupsen/logrus"
)

func (c *Converter) GetInfoURLFromDB() (map[string]string, error) {
	urls := []*model.DescriptionURL{}
	if err := c.db.Select(&urls, "SELECT uuid as id, short_url as short, original_url as original FROM urls.links"); err != nil {
		return nil, err
	}
	result := c.convertDescriptionURLToMap(urls)
	return result, nil
}

func (c *Converter) StoreURLInDB(shortURL string, originalURL string) error {
	result, err := c.db.Exec(`INSERT INTO urls.links (short_url, original_url)
								VALUES ($1, $2)
								ON CONFLICT (short_url) DO NOTHING;
							 `, shortURL, originalURL)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	logrus.Printf("Добавлено в БД: %d строк", rowsAffected)
	return nil
}
