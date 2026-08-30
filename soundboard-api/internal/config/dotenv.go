package config

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// DefaultDotEnvName is the file LoadDotEnv looks for when walking up from the working
// directory.
const DefaultDotEnvName = ".env"

// LoadDotEnv reads key=value pairs from path into the process environment.
//
// This exists for the CLI, not the API. The service is handed its environment by
// Railway and should never read a file for it, but the CLI is run by hand from a
// PowerShell window that starts with nothing set — and because DBPath and AudioDir both
// have working local defaults, an unset environment does not fail loudly. It silently
// writes to the local database and the local clips folder instead of Turso and R2, which
// looks exactly like a successful publish.
//
// A variable already present in the real environment always wins, so an explicit
//
//	$env:TURSO_DATABASE_URL = "..."
//
// still overrides the file for a one-off run.
//
// A missing file is not an error: a checkout with no credentials is a supported way to
// work, and it is what the local defaults are for.
func LoadDotEnv(path string) error {
	file, err := os.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for line := 1; scanner.Scan(); line++ {
		key, value, ok := parseDotEnvLine(scanner.Text())
		if !ok {
			continue
		}
		if key == "" {
			return fmt.Errorf("%s:%d: name is empty", path, line)
		}
		// Setenv rather than Setenv-if-absent: the real environment is checked first so
		// a value exported in the shell is not clobbered by a stale file.
		if _, set := os.LookupEnv(key); set {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("%s:%d: set %s: %w", path, line, key, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	return nil
}

// parseDotEnvLine splits one line into a name and value, reporting whether the line
// carried an assignment at all. Blank lines and comments return ok=false.
//
// Leading whitespace is tolerated on both the name and the value, because a hand-edited
// file picks up indentation easily and silently ignoring an indented line is the worst
// possible response to it. An optional `export ` prefix is accepted so a file written
// for bash also works here.
func parseDotEnvLine(raw string) (key, value string, ok bool) {
	// Trim \r as well as spaces: a file saved by a Windows editor keeps CRLF endings,
	// and a trailing \r would otherwise become part of the value.
	line := strings.TrimSpace(strings.TrimSuffix(raw, "\r"))
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}

	line = strings.TrimPrefix(line, "export ")

	name, val, found := strings.Cut(line, "=")
	if !found {
		return "", "", false
	}
	return strings.TrimSpace(name), unquote(strings.TrimSpace(val)), true
}

// unquote strips one matching pair of surrounding quotes. Quoting is what lets a value
// keep leading or trailing spaces that TrimSpace would otherwise eat.
func unquote(value string) string {
	if len(value) < 2 {
		return value
	}
	first, last := value[0], value[len(value)-1]
	if first != last {
		return value
	}
	if first == '"' || first == '\'' {
		return value[1 : len(value)-1]
	}
	return value
}

// FindDotEnv looks for a .env beside the working directory, then in each parent up to
// limit levels. The CLI is run from soundboard-api/ most of the time but from the repo
// root often enough that requiring one exact location is a papercut.
//
// It returns an empty string when nothing is found, which LoadDotEnv treats as fine.
func FindDotEnv(startDir string, limit int) string {
	dir := startDir
	for i := 0; i <= limit; i++ {
		candidate := filepath.Join(dir, DefaultDotEnvName)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
