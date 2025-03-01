package turso

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Repository provides access to the Turso database
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new Turso repository instance
func NewRepository(dbURL, authToken string) (*Repository, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	// Add authentication token to the URL if provided
	if authToken != "" {
		dbURL = fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
	}

	// Connect to the database
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Turso database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to ping Turso database: %w (close error: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to ping Turso database: %w", err)
	}

	return &Repository{
		db: db,
	}, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// Exec executes a query without returning any rows
func (r *Repository) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

// Query executes a query that returns rows
func (r *Repository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that returns a single row
func (r *Repository) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.db.QueryRowContext(ctx, query, args...)
}
