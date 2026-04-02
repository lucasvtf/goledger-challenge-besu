package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var ErrNoSyncedValue = errors.New("no synced value found")

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS contract_state (
			id SERIAL PRIMARY KEY,
			value BIGINT NOT NULL,
			synced_at TIMESTAMP DEFAULT NOW()
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error running migration: %w", err)
	}
	return nil
}

func SaveValue(db *sql.DB, value int64) error {
	_, err := db.Exec("INSERT INTO contract_state (value) VALUES ($1)", value)
	if err != nil {
		return fmt.Errorf("error saving value: %w", err)
	}
	return nil
}

func GetLatestValue(db *sql.DB) (int64, error) {
	var value int64
	err := db.QueryRow("SELECT value FROM contract_state ORDER BY id DESC LIMIT 1").Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNoSyncedValue
	}
	if err != nil {
		return 0, fmt.Errorf("error querying latest value: %w", err)
	}
	return value, nil
}
