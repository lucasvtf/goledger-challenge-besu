CREATE TABLE IF NOT EXISTS contract_values (
    id SERIAL PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO contract_values (value) VALUES ('0');