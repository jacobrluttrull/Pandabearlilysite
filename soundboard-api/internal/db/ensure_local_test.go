package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A fresh checkout has no data/ directory, so this is the case that has to work: the
// parent is created and opening the file afterwards succeeds.
func TestEnsureLocalDirCreatesMissingParent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "soundboard.db")

	if err := EnsureLocalDir(path); err != nil {
		t.Fatalf("EnsureLocalDir() = %v, want nil", err)
	}

	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open() after EnsureLocalDir = %v, want nil", err)
	}
	db.Close()
}

// Called twice — once by the API and once by the CLI against the same path — it must not
// object to a directory that is already there.
func TestEnsureLocalDirAcceptsExistingParent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "soundboard.db")

	if err := EnsureLocalDir(path); err != nil {
		t.Fatalf("first call = %v, want nil", err)
	}
	if err := EnsureLocalDir(path); err != nil {
		t.Fatalf("second call = %v, want nil", err)
	}
}

// The deployment case: the parent cannot be created, and the message has to point at the
// missing Turso configuration rather than at the filesystem, because that is the actual
// mistake. A regular file stands in for the container's root-owned /app.
func TestEnsureLocalDirNamesTursoWhenParentUnusable(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "data")
	if err := writeFile(blocker); err != nil {
		t.Fatalf("setup: %v", err)
	}

	err := EnsureLocalDir(filepath.Join(blocker, "soundboard.db"))
	if err == nil {
		t.Fatal("EnsureLocalDir() = nil, want an error when the parent cannot be created")
	}
	if !strings.Contains(err.Error(), "TURSO_DATABASE_URL") {
		t.Errorf("error does not mention TURSO_DATABASE_URL: %v", err)
	}
}

// writeFile creates an empty regular file at path.
func writeFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}
