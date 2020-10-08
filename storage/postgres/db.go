package postgres

import (
	"database/sql"
	"fmt"
)

func NewDBFromURL(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("opening postgres: %w", err)
	}
	return db, nil
}
