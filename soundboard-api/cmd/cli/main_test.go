package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvFileDoesNothingWithoutExplicitSelection(t *testing.T) {
	t.Setenv("SOUNDBOARD_ENV_FILE", "")
	if err := loadEnvFile(); err != nil {
		t.Fatalf("loadEnvFile() with no selected file: %v", err)
	}
}

func TestLoadEnvFileReadsExplicitlySelectedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "production.env")
	if err := os.WriteFile(path, []byte("R2_BUCKET=production-clips\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SOUNDBOARD_ENV_FILE", path)
	t.Setenv("R2_BUCKET", "")
	os.Unsetenv("R2_BUCKET")

	if err := loadEnvFile(); err != nil {
		t.Fatalf("loadEnvFile(): %v", err)
	}
	if got := os.Getenv("R2_BUCKET"); got != "production-clips" {
		t.Errorf("R2_BUCKET = %q, want production-clips", got)
	}
}

func TestLoadEnvFileRejectsMissingExplicitFile(t *testing.T) {
	t.Setenv("SOUNDBOARD_ENV_FILE", filepath.Join(t.TempDir(), "missing.env"))
	if err := loadEnvFile(); err == nil {
		t.Error("loadEnvFile() succeeded for a missing selected file")
	}
}
