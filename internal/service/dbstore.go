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

func (c *Converter) TranStoreURLInDB(urls []*model.DescriptionURL, user string) error {
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
									RETURNING short_url, uuid;
								`)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	insertUserStmt, err := tx.Preparex(`
        INSERT INTO urls.users (id) VALUES ($1)
        ON CONFLICT (id) DO NOTHING;
    `)
	if err != nil {
		return err
	}
	defer insertUserStmt.Close()

	insertXMapStmt, err := tx.Preparex(`
        INSERT INTO urls.links_users_xmap (user_id, url_id)
        VALUES ($1, $2)
        ON CONFLICT (url_id) DO NOTHING;
    `)
	if err != nil {
		return err
	}
	defer insertXMapStmt.Close()

	if user != "" {
		_, err = insertUserStmt.Exec(user)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	for _, url := range urls {
		var linkID string

		if err := insertStmt.QueryRow(url.Short, url.Original).Scan(&url.Short, &linkID); err != nil {
			_ = tx.Rollback()
			return err
		}

		if user != "" {
			_, err = insertXMapStmt.Exec(user, linkID)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	logrus.Printf("Данные в БД успешно обновлены")
	return nil
}

func (c *Converter) GetInfoUserURLFromDB(user string) ([]*model.DescriptionURL, error) {
	var result []*model.DescriptionURL

	query := `
        SELECT  l.short_url, l.original_url
        FROM urls.links_users_xmap x
        INNER JOIN urls.links l ON l.uuid = x.url_id 
		WHERE x.user_id = $1
    `

	if err := c.db.Select(&result, query, user); err != nil {
		return nil, err
	}

	return result, nil
}
