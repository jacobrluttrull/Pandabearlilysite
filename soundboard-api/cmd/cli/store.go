package main

import (
	"database/sql"
	"fmt"

	"soundboard-api/internal/config"
	sbdb "soundboard-api/internal/db"
	"soundboard-api/internal/db/gen"
)

// store bundles the handles every command needs: the open database, the generated query
// set, and the resolved config.
type store struct {
	db      *sql.DB
	queries *gen.Queries
	cfg     config.Config
}

// openStore loads config, opens the database, and brings it up to the latest schema.
//
// Every command starts this way, so routing it through one function means the migration
// step cannot be forgotten when a new command is added.
func openStore() (*store, error) {
	cfg := config.Load()

	sqlDB, err := sbdb.Open(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := sbdb.Migrate(sqlDB); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return &store{db: sqlDB, queries: gen.New(sqlDB), cfg: cfg}, nil
}

func (s *store) Close() error {
	return s.db.Close()
}
