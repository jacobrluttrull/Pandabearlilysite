package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Open opens the SQLite database at path, creating the file if it doesn't exist.
func Open(path string) (*sql.DB, error) {
	// WAL matters now that the API writes on every play, not just when the CLI uploads:
	// readers no longer block behind the current writer, and each write costs one
	// append instead of a rollback-journal round trip. busy_timeout makes a contended
	// write wait briefly rather than failing outright with SQLITE_BUSY.
	dsn := path + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	// SQLite only supports one writer at a time; a single connection avoids
	// SQLITE_BUSY errors under concurrent access. This also serialises reads, which is
	// fine at fan-site traffic — if plays ever outgrow it, the next step is separate
	// read and write pools rather than a bigger single pool.
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite db: %w", err)
	}

	return db, nil
}
