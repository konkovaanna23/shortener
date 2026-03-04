package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/sirupsen/logrus"
)

// GetOriginalURLFromDB получает оригинальный URL по короткому из БД.
func (c *Converter) GetOriginalURLFromDB(shortURL string) (string, error) {
	var resultOriginal string
	var resultIsDeleted bool
	err := c.db.QueryRow(` SELECT original_url, is_deleted
						   FROM urls.links 
						   WHERE short_url=$1; `, shortURL).Scan(&resultOriginal, &resultIsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("не существует оригинально URL для %s", shortURL)
		}
		return "", err
	}
	if resultIsDeleted {
		return "", model.ErrorDeletedURL
	}
	return resultOriginal, nil
}

// TranStoreURLInDB сохраняет URL в БД.
func (c *Converter) TranStoreURLInDB(urls []*model.DescriptionURL, user string) error {
	var errConflict error
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
		var urlShort string

		err = insertStmt.QueryRow(url.Short, url.Original).Scan(&urlShort, &linkID)

		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if urlShort != url.Short {
			url.Short = urlShort
			errConflict = model.ErrorConflictURL
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
	return errConflict
}

// GetInfoUserURLFromDB получает информацию URL по пользователю.
func (c *Converter) GetInfoUserURLFromDB(user string) ([]*model.DescriptionURL, error) {
	var result []*model.DescriptionURL

	query := `
        SELECT  l.short_url, l.original_url
        FROM urls.links_users_xmap x
        INNER JOIN urls.links l ON l.uuid = x.url_id 
		WHERE x.user_id = $1 and l.is_deleted = false
    `

	if err := c.db.Select(&result, query, user); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Converter) GetStatsDB() (*model.Stats, error) {
	var stats model.Stats
	query := `
        SELECT 
            (SELECT count(*) FROM urls.links) AS urls,  
            (SELECT count(*) FROM urls.users) AS users
    `

	if err := c.db.Get(&stats, query); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (c *Converter) TranDeleteURLsFromDB(urls []*model.DescriptionURL) error {
	if len(urls) == 0 {
		logrus.Warningln("Нет данных для удаления")
		return nil
	}
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Preparex(`
		UPDATE urls.links 
		SET is_deleted = true 
		WHERE short_url = $1 
		  AND uuid IN (
		    SELECT url_id 
		    FROM urls.links_users_xmap 
		    WHERE user_id = $2
		  )
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, url := range urls {
		if url.Short == "" || url.UserID == "" {
			return fmt.Errorf("не задан url или user")
		}

		result, err := stmt.Exec(url.Short, url.UserID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if rowsAffected == 0 {
			logrus.Infof("URL=%s не существует или не принадлежит пользователю=%s", url.Short, url.UserID)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	logrus.Info("URL успешны удалены")
	return nil
}
