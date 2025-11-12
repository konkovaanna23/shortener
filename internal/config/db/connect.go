package db

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewConnect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := Ping(db); err != nil {
		return nil, err
	}

	return db, nil
}

func Ping(db *sqlx.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}
