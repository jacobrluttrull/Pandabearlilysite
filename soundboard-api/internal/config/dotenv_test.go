package config

import (
	"os"
	"path/filepath"
	"testing"
)

// writeEnv drops a .env with the given contents into a temp dir and returns its path.
func writeEnv(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func TestLoadDotEnvReadsValues(t *testing.T) {
	// Indented lines and CRLF endings are both what a hand-edited file on Windows
	// actually looks like — the real .env in this repo was indented with two spaces.
	path := writeEnv(t, "# a comment\r\n"+
		"  R2_BUCKET=pandalily-clips\r\n"+
		"\r\n"+
		"export TURSO_DATABASE_URL=libsql://example.turso.io\r\n"+
		"QUOTED=\"  padded  \"\r\n"+
		"EMPTY=\r\n")

	for _, key := range []string{"R2_BUCKET", "TURSO_DATABASE_URL", "QUOTED", "EMPTY"} {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}

	if err := LoadDotEnv(path); err != nil {
		t.Fatalf("LoadDotEnv: %v", err)
	}

	cases := map[string]string{
		"R2_BUCKET":          "pandalily-clips",
		"TURSO_DATABASE_URL": "libsql://example.turso.io",
		"QUOTED":             "  padded  ",
		"EMPTY":              "",
	}
	for key, want := range cases {
		if got := os.Getenv(key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}

// The real environment is the more deliberate signal: a value exported in the shell for
// a one-off run must not be silently replaced by a stale file.
func TestLoadDotEnvDoesNotOverrideRealEnvironment(t *testing.T) {
	path := writeEnv(t, "R2_BUCKET=from-file\n")

	t.Setenv("R2_BUCKET", "from-shell")
	if err := LoadDotEnv(path); err != nil {
		t.Fatalf("LoadDotEnv: %v", err)
	}
	if got := os.Getenv("R2_BUCKET"); got != "from-shell" {
		t.Errorf("R2_BUCKET = %q, want the shell value to win", got)
	}
}

// A checkout with no credentials is a supported way to work, so an absent file is not a
// failure — the local defaults are exactly what it should fall back to.
func TestLoadDotEnvMissingFileIsNotAnError(t *testing.T) {
	if err := LoadDotEnv(filepath.Join(t.TempDir(), "nothing-here")); err != nil {
		t.Errorf("missing .env should be tolerated, got %v", err)
	}
}

func TestFindDotEnvWalksUp(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "soundboard-api", "cmd", "cli")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	want := filepath.Join(root, ".env")
	if err := os.WriteFile(want, []byte("K=v\n"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	if got := FindDotEnv(nested, 3); got != want {
		t.Errorf("FindDotEnv = %q, want %q", got, want)
	}
	// Out of reach: the limit stops the walk before the file.
	if got := FindDotEnv(nested, 1); got != "" {
		t.Errorf("FindDotEnv beyond limit = %q, want empty", got)
	}
}
