package config

import (
	"net/url"
	"testing"
)

func TestDatabaseDSN(t *testing.T) {
	cases := []struct {
		name, dbPath, databaseURL, authToken, want string
	}{
		{
			"falls back to the local file when no URL is set",
			"data/soundboard.db", "", "",
			"data/soundboard.db",
		},
		{
			"ignores a stray token when there is no URL",
			"data/soundboard.db", "", "leftover-token",
			"data/soundboard.db",
		},
		{
			"passes the URL through untouched when there is no token",
			"data/soundboard.db", "libsql://soundboard-user.turso.io", "",
			"libsql://soundboard-user.turso.io",
		},
		{
			"appends the token as a query parameter",
			"data/soundboard.db", "libsql://soundboard-user.turso.io", "tok123",
			"libsql://soundboard-user.turso.io?authToken=tok123",
		},
		{
			"keeps an existing query parameter",
			"data/soundboard.db", "libsql://soundboard-user.turso.io?syncInterval=60", "tok123",
			"libsql://soundboard-user.turso.io?authToken=tok123&syncInterval=60",
		},
		{
			"escapes a token containing url-significant characters",
			"data/soundboard.db", "libsql://soundboard-user.turso.io", "a b&c=d/e+f",
			"libsql://soundboard-user.turso.io?authToken=a+b%26c%3Dd%2Fe%2Bf",
		},
		{
			"replaces a token already present on the URL",
			"data/soundboard.db", "libsql://soundboard-user.turso.io?authToken=stale", "fresh",
			"libsql://soundboard-user.turso.io?authToken=fresh",
		},
		{
			"returns a malformed URL unchanged instead of panicking",
			"data/soundboard.db", "libsql://soundboard-user.turso.io/%zz", "tok123",
			"libsql://soundboard-user.turso.io/%zz",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := Config{DBPath: c.dbPath, DatabaseURL: c.databaseURL, AuthToken: c.authToken}
			if got := cfg.DatabaseDSN(); got != c.want {
				t.Errorf("DatabaseDSN() = %q, want %q", got, c.want)
			}
		})
	}
}

// The escaped form only matters if the driver can read the original token back out, so
// check the round trip rather than trusting the encoded spelling alone.
func TestDatabaseDSNTokenRoundTrips(t *testing.T) {
	for _, token := range []string{
		"tok123", "a b&c=d/e+f", "eyJhbGciOiJFZERTQSJ9.payload-part.signature?&=", "üñí",
	} {
		cfg := Config{
			DBPath:      "data/soundboard.db",
			DatabaseURL: "libsql://soundboard-user.turso.io?syncInterval=60",
			AuthToken:   token,
		}
		u, err := url.Parse(cfg.DatabaseDSN())
		if err != nil {
			t.Fatalf("DatabaseDSN() for token %q did not parse: %v", token, err)
		}
		if got := u.Query().Get("authToken"); got != token {
			t.Errorf("authToken = %q, want %q", got, token)
		}
		if got := u.Query().Get("syncInterval"); got != "60" {
			t.Errorf("syncInterval = %q, want %q", got, "60")
		}
	}
}
