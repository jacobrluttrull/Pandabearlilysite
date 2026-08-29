package db

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureLocalDir creates the directory holding a local SQLite file, and explains the
// failure in terms of the choice that led to it when it cannot.
//
// Two situations reach this, and the raw driver error serves neither. A fresh checkout
// has no data/ directory, because it is gitignored — SQLite will not create the parent of
// a database file, so cloning and running failed with "unable to open database file (14)"
// and no hint that a mkdir was all it needed. Creating it here is what makes the local
// default actually zero-setup.
//
// The other is a deployment that never meant to be here at all. In the container /app is
// root-owned and the service runs as an unprivileged user, so the mkdir fails, and the
// real mistake is upstream: TURSO_DATABASE_URL was missing or misspelled, the local file
// was a silent fallback, and the resulting error named neither Turso nor the fallback.
// Saying so costs one sentence and saves reading it as a bucket or migration problem.
func EnsureLocalDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf(
			"cannot create local database directory %q: %w.\n"+
				"If this is a deployment, that is the real problem: TURSO_DATABASE_URL is unset or "+
				"misspelled, so the database silently fell back to a local file. Set TURSO_DATABASE_URL "+
				"and TURSO_AUTH_TOKEN. If this is local development, point SOUNDBOARD_DB_PATH at a "+
				"writable location",
			dir, err)
	}
	return nil
}
