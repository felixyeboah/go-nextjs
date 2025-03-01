package database

import (
	"database/sql"
	"fmt"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type TursoDB struct {
	*sql.DB
}

func NewTursoConnection(url, authToken string) (*TursoDB, error) {
	if url == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	// Add authentication token to the URL if provided
	if authToken != "" {
		url = fmt.Sprintf("%s?authToken=%s", url, authToken)
	}

	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return &TursoDB{db}, nil
}

func (db *TursoDB) Close() error {
	return db.DB.Close()
} 