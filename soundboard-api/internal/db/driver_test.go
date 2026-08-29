package db

import (
	"strings"
	"testing"
)

func TestDriverFor(t *testing.T) {
	cases := []struct {
		name, dsn, wantDriver, wantDSN string
	}{
		{
			"turso url goes remote untouched",
			"libsql://soundboard-pandalily.turso.io?authToken=abc123",
			driverLibSQL,
			"libsql://soundboard-pandalily.turso.io?authToken=abc123",
		},
		{
			"websocket url goes remote",
			"wss://soundboard-pandalily.turso.io?authToken=abc123",
			driverLibSQL,
			"wss://soundboard-pandalily.turso.io?authToken=abc123",
		},
		{
			"https url goes remote",
			"https://soundboard-pandalily.turso.io?authToken=abc123",
			driverLibSQL,
			"https://soundboard-pandalily.turso.io?authToken=abc123",
		},
		{
			"relative path goes local with pragmas",
			"data/soundboard.db",
			driverSQLite,
			"data/soundboard.db" + localPragmas,
		},
		{
			"windows absolute path is a path, not a scheme",
			`C:\data\soundboard.db`,
			driverSQLite,
			`C:\data\soundboard.db` + localPragmas,
		},
		{
			"unix absolute path goes local",
			"/data/soundboard.db",
			driverSQLite,
			"/data/soundboard.db" + localPragmas,
		},
		{
			"in-memory database goes local",
			":memory:",
			driverSQLite,
			":memory:" + localPragmas,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotDriver, gotDSN := driverFor(c.dsn)
			if gotDriver != c.wantDriver {
				t.Errorf("driverFor(%q) driver = %q, want %q", c.dsn, gotDriver, c.wantDriver)
			}
			if gotDSN != c.wantDSN {
				t.Errorf("driverFor(%q) dsn = %q, want %q", c.dsn, gotDSN, c.wantDSN)
			}
		})
	}
}

func TestIsRemoteDSN(t *testing.T) {
	cases := []struct {
		name string
		dsn  string
		want bool
	}{
		{"libsql scheme", "libsql://db.turso.io", true},
		{"wss scheme", "wss://db.turso.io", true},
		{"https scheme", "https://db.turso.io", true},
		{"relative path", "data/soundboard.db", false},
		{"windows drive letter", `C:\data\soundboard.db`, false},
		{"windows path with forward slashes", "C:/data/soundboard.db", false},
		{"unix path", "/var/lib/soundboard/soundboard.db", false},
		{"empty dsn", "", false},
		{"scheme name without separator", "libsql:soundboard.db", false},
		{"ws scheme", "ws://127.0.0.1:8080", true},
		{"http scheme, as served by `turso dev`", "http://127.0.0.1:8080", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsRemoteDSN(c.dsn); got != c.want {
				t.Errorf("IsRemoteDSN(%q) = %v, want %v", c.dsn, got, c.want)
			}
		})
	}
}

// Pragmas configure how SQLite manages a file on disk, so a remote DSN must never carry
// them — the query string it already has is the auth token, and appending to it would
// both corrupt that and send file settings to a server with no file.
func TestRemoteDSNCarriesNoPragmas(t *testing.T) {
	const dsn = "libsql://soundboard-pandalily.turso.io?authToken=abc123"

	_, gotDSN := driverFor(dsn)
	if strings.Contains(gotDSN, "_pragma") {
		t.Errorf("driverFor(%q) dsn = %q, want no pragmas", dsn, gotDSN)
	}
}

// RequireRemote is the guard against a mistyped TURSO_DATABASE_URL silently becoming a
// filesystem path. The bare-hostname case is the one seen in practice: Turso's dashboard
// can show the host without its scheme, and copying that alone parses as a relative path.
func TestRequireRemote(t *testing.T) {
	cases := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{"libsql url", "libsql://db.turso.io?authToken=x", false},
		{"turso dev http url", "http://127.0.0.1:8080", false},
		{"bare hostname, scheme omitted", "soundboard-you.turso.io", true},
		{"relative file path", "data/soundboard.db", true},
		{"empty", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := RequireRemote(c.dsn)
			if (err != nil) != c.wantErr {
				t.Errorf("RequireRemote(%q) error = %v, wantErr %v", c.dsn, err, c.wantErr)
			}
		})
	}
}
