package config

import "net/url"

// DatabaseDSN returns the DSN to hand to db.Open.
//
// The auth token rides inside the DSN rather than being passed alongside it because
// database/sql only ever gives a driver one string: sql.Open(driverName, dsn). There is
// no second argument to put a credential in, so libSQL — like Postgres and MySQL before
// it — expects the credential as a query parameter on the URL. That makes the returned
// value secret whenever AuthToken is set: it must not be logged or wrapped into an error.
//
// With no DatabaseURL there is no remote and no credential, so the local file path is
// the DSN as-is; that is the default for local development.
func (c Config) DatabaseDSN() string {
	if c.DatabaseURL == "" {
		return c.DBPath
	}
	if c.AuthToken == "" {
		return c.DatabaseURL
	}

	// Parse rather than concatenate "?authToken=": the URL may already carry a query
	// (Turso hands out URLs with options on them), and a token can contain characters
	// like "=" or "." that need escaping to survive as a single parameter value.
	u, err := url.Parse(c.DatabaseURL)
	if err != nil {
		// A URL this broken cannot be repaired here, and config load is the wrong place
		// to fail: returning it untouched lets the driver report a connection error with
		// its own context instead of the process dying at startup with less to say.
		return c.DatabaseURL
	}
	q := u.Query()
	q.Set("authToken", c.AuthToken)
	u.RawQuery = q.Encode()
	return u.String()
}
