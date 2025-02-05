package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DB struct {
	Pool *pgxpool.Pool
}

// NewDB loads environment variables from .env and initializes the database connection
func NewDB() (*DB, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Not fatal if .env doesn't exist; environment variables may be set elsewhere
		fmt.Printf("Warning: .env file not found or could not be loaded: %v\n", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// SaveProofData inserts the username, proof, and helper into the database
func (db *DB) SaveProofData(ctx context.Context, username string, proof []byte, helper_data []byte) error {
	query := `INSERT INTO zetable (username, proof, helper_data) VALUES ($1, $2, $3)`
	_, err := db.Pool.Exec(ctx, query, username, proof, helper_data)
	if err != nil {
		return fmt.Errorf("failed to insert proof data: %w", err)
	}
	return nil
}

// GetProofData retrieves the proof and helper_data for a given username.
// Returns proofBytes, helperBytes, and an error if any.
func (db *DB) GetProofData(ctx context.Context, username string) ([]byte, []byte, error) {
	query := `SELECT proof, helper_data FROM zetable WHERE username = $1 LIMIT 1`
	row := db.Pool.QueryRow(ctx, query, username)

	var proofBytes, helperBytes []byte
	err := row.Scan(&proofBytes, &helperBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve proof data for username %q: %w", username, err)
	}

	return proofBytes, helperBytes, nil
}

// Close closes the database connection pool.
func (db *DB) Close() {
	db.Pool.Close()
}
