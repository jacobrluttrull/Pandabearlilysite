package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Open opens the SQLite database at path, creating the file if it doesn't exist.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	// SQLite only supports one writer at a time; a single connection avoids
	// SQLITE_BUSY errors under concurrent access.
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite db: %w", err)
	}

	return db, nil
}
