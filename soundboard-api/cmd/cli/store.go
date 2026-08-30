package main

import (
	"database/sql"
	"fmt"
	"os"

	"soundboard-api/internal/clipstore"
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
	// clips is where audio is published to. Commands still read and hash the local
	// files directly — that is where the source of truth for a new clip is — but
	// anything that changes the stored collection goes through here, so the same
	// command works whether it is writing to a folder or to R2.
	clips clipstore.Manager
}

// openStore loads config, opens the database, and brings it up to the latest schema.
//
// Every command starts this way, so routing it through one function means the migration
// step cannot be forgotten when a new command is added.
func openStore() (*store, error) {
	cfg := config.Load()

	// TURSO_DATABASE_URL being set states an intent to use a remote database. If the
	// value is malformed it would otherwise be taken as a filesystem path and the
	// service would start quietly on a local file that no deploy survives.
	if cfg.DatabaseURL != "" {
		if err := sbdb.RequireRemote(cfg.DatabaseDSN()); err != nil {
			return nil, err
		}
	} else if err := sbdb.EnsureLocalDir(cfg.DBPath); err != nil {
		return nil, err
	}

	sqlDB, err := sbdb.Open(cfg.DatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := sbdb.Migrate(sqlDB); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	if err := clipstore.RequireComplete(cfg.R2()); err != nil {
		sqlDB.Close()
		return nil, err
	}

	clips, ok := clipstore.Open(cfg.R2(), cfg.AudioDir).(clipstore.Manager)
	if !ok {
		sqlDB.Close()
		return nil, fmt.Errorf("clip store is read-only, cannot manage clips")
	}

	// Printed on every command, to stderr so it never contaminates piped output. The
	// failure this guards against is not an error but a success against the wrong
	// target: with no environment set, DBPath and AudioDir both have working local
	// defaults, so an import writes to a file on this machine and reports that it
	// worked. Naming the target makes that visible in the one place it matters.
	fmt.Fprintf(os.Stderr, "using %s, clips in %s\n\n", describeDatabase(cfg), clips.Describe())

	return &store{db: sqlDB, queries: gen.New(sqlDB), cfg: cfg, clips: clips}, nil
}

// describeDatabase names the database being used without revealing how to reach it.
// DatabaseURL is a hostname and safe to show; the auth token that accompanies it is a
// credential and never appears here, which is why this does not print the DSN.
func describeDatabase(cfg config.Config) string {
	if cfg.DatabaseURL != "" {
		return "remote database " + cfg.DatabaseURL
	}
	return "local database file " + cfg.DBPath
}

func (s *store) Close() error {
	return s.db.Close()
}
