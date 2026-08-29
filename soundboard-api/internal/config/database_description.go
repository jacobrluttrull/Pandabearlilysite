package config

import "net/url"

// DatabaseDescription names the database in a form safe to log.
//
// This exists because the fallback in DatabaseDSN is silent by design: with no
// TURSO_DATABASE_URL set, the service happily opens a local file and runs. That is right
// for development and dangerous in production, where the container filesystem does not
// survive a deploy — play counts would accumulate into a file that is thrown away, with
// nothing in the logs to say so. Printing this at startup makes the choice visible in
// Railway's log the moment the process boots, rather than at the first lost tally.
//
// The auth token is never included: only the host is taken from the URL, so the returned
// string carries no credential and can be logged freely.
func (c Config) DatabaseDescription() string {
	if c.DatabaseURL == "" {
		return "local file " + c.DBPath
	}

	u, err := url.Parse(c.DatabaseURL)
	if err != nil || u.Host == "" {
		// Unparseable, so there is no host to name. Say which database was configured
		// without echoing a string that may contain anything.
		return "remote libSQL (url unparseable)"
	}
	return "remote libSQL at " + u.Scheme + "://" + u.Host
}
