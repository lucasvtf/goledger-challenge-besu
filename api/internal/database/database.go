package database

import (
	"database/sql"
	"fmt"
	"time"

	"besu-api/internal/config"

	_ "github.com/lib/pq"
)

type Client struct {
	db *sql.DB
}

type ContractValue struct {
	ID        int       `json:"id"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewClient(cfg *config.Config) (*Client, error) {
	var connectionString string
	
	if cfg.DatabaseURL != "" {
		connectionString = cfg.DatabaseURL
	} else {
		connectionString = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DatabaseHost,
			cfg.DatabasePort,
			cfg.DatabaseUser,
			cfg.DatabasePass,
			cfg.DatabaseName,
		)
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	client := &Client{db: db}

	if err := client.createTable(); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return client, nil
}

func (c *Client) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS contract_values (
		id SERIAL PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	INSERT INTO contract_values (value) 
	SELECT '0' 
	WHERE NOT EXISTS (SELECT 1 FROM contract_values);
	`

	_, err := c.db.Exec(query)
	return err
}

func (c *Client) GetValue() (*ContractValue, error) {
	query := `
	SELECT id, value, updated_at 
	FROM contract_values 
	ORDER BY updated_at DESC 
	LIMIT 1
	`

	var cv ContractValue
	err := c.db.QueryRow(query).Scan(&cv.ID, &cv.Value, &cv.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SetValue("0")
		}
		return nil, fmt.Errorf("failed to get value from database: %v", err)
	}

	return &cv, nil
}

func (c *Client) SetValue(value string) (*ContractValue, error) {
	updateQuery := `
	UPDATE contract_values 
	SET value = $1, updated_at = CURRENT_TIMESTAMP 
	WHERE id = (SELECT id FROM contract_values ORDER BY updated_at DESC LIMIT 1)
	RETURNING id, value, updated_at
	`

	var cv ContractValue
	err := c.db.QueryRow(updateQuery, value).Scan(&cv.ID, &cv.Value, &cv.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			insertQuery := `
			INSERT INTO contract_values (value, updated_at) 
			VALUES ($1, CURRENT_TIMESTAMP) 
			RETURNING id, value, updated_at
			`
			
			err = c.db.QueryRow(insertQuery, value).Scan(&cv.ID, &cv.Value, &cv.UpdatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to insert initial value in database: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to update value in database: %v", err)
		}
	}

	return &cv, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}