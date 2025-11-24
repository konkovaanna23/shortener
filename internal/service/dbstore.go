package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/sirupsen/logrus"
)

func (c *Converter) GetOriginalURLFromDB(shortURL string) (string, error) {
	fmt.Println(shortURL)
	var resultOriginal string
	err := c.db.QueryRow(` SELECT original_url 
						   FROM urls.links 
						   WHERE short_url=$1; `, shortURL).Scan(&resultOriginal)
	fmt.Println(err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("не существует оригинально URL для %s", shortURL)
		}
		return "", err
	}
	return resultOriginal, nil
}

func (c *Converter) StoreURLInDB(shortURL, originalURL string) (string, error) {
	var resultShort string

	err := c.db.QueryRow(`
        INSERT INTO urls.links (short_url, original_url)
        VALUES ($1, $2)
        ON CONFLICT (original_url) DO UPDATE
		SET short_url = urls.links.short_url
        RETURNING short_url;
    `, shortURL, originalURL).Scan(&resultShort)

	if err != nil {
		return "", err
	}

	if shortURL != resultShort {
		return resultShort, model.ErrorConflictURL
	}
	return resultShort, nil
}

func (c *Converter) TranStoreURLInDB(urls []*model.DescriptionURL) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	insertStmt, err := tx.Prepare(`INSERT INTO urls.links (short_url, original_url)
									VALUES ($1, $2)
									ON CONFLICT (original_url) DO UPDATE
									SET short_url = urls.links.short_url
									RETURNING short_url;
								`)
	if err != nil {
		return err
	}
	defer insertStmt.Close()
	for _, url := range urls {
		if err := insertStmt.QueryRow(url.Short, url.Original).Scan(&url.Short); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	logrus.Printf("Данные в БД обновлены")
	return nil
}
