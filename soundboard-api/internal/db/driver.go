package db

import (
	"fmt"
	"strings"
)

// The two database drivers this binary can talk to. Both are pure Go, which is what
// keeps CGO_ENABLED=0 builds working — the container image is built static, so a driver
// needing cgo (go-libsql, mattn/go-sqlite3) would break the Dockerfile, not just slow
// the build down.
const (
	driverSQLite = "sqlite" // modernc.org/sqlite, a local file on disk
	driverLibSQL = "libsql" // libsql-client-go, a Turso database over the network
)

// localPragmas are appended to a local file DSN. These are all statements about how
// SQLite manages a file on disk, so they are meaningless — and rejected — against a
// remote database, where the server owns its own journalling and locking.
//
// WAL matters now that the API writes on every play, not just when the CLI uploads:
// readers no longer block behind the current writer, and each write costs one append
// instead of a rollback-journal round trip. busy_timeout makes a contended write wait
// briefly rather than failing outright with SQLITE_BUSY.
const localPragmas = "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"

// remoteSchemes mark a DSN as pointing at a libSQL server rather than a file.
//
// libsql:// is what Turso Cloud hands out; the rest are the same endpoint named by the
// transport it rides on, and the client accepts all of them. The plaintext pair matters
// for development, not production: `turso dev` runs a libSQL server on localhost over
// http://, which is how the remote code path gets exercised without cloud credentials.
//
// Matching on the full "scheme://" rather than just the scheme name is deliberate: a
// Windows path like C:\data\soundboard.db has a colon in position two and would
// otherwise look like a URL scheme to any parser loose enough to accept it.
var remoteSchemes = []string{"libsql://", "wss://", "ws://", "https://", "http://"}

// IsRemoteDSN reports whether dsn addresses a network database instead of a local file.
//
// Exported because the entrypoints need it to catch a mistyped TURSO_DATABASE_URL.
// Setting that variable states an intent to use a remote database; if the value does not
// parse as one, the DSN would otherwise be taken as a filesystem path and the service
// would start quietly on a local file. See RequireRemote.
func IsRemoteDSN(dsn string) bool {
	for _, scheme := range remoteSchemes {
		if strings.HasPrefix(dsn, scheme) {
			return true
		}
	}
	return false
}

// driverFor picks the driver for a DSN and returns the connection string to hand it.
//
// A remote DSN passes through untouched — it already carries its ?authToken=... query,
// and appending anything would either corrupt that or send file pragmas to a server that
// has no file. Anything else is taken as a filesystem path and gets the local pragmas.
func driverFor(dsn string) (driverName, connDSN string) {
	if IsRemoteDSN(dsn) {
		return driverLibSQL, dsn
	}
	return driverSQLite, dsn + localPragmas
}

// RequireRemote returns an error unless dsn addresses a remote database.
//
// Call it when configuration says a remote database was intended. Falling back to a
// local file there is the worst possible outcome: on a container host the file is
// discarded at the next deploy, so the service appears healthy while quietly losing
// every write. Refusing to start is the safer failure.
func RequireRemote(dsn string) error {
	if IsRemoteDSN(dsn) {
		return nil
	}
	return fmt.Errorf(
		"TURSO_DATABASE_URL is set but is not a remote database URL: it must start with one of %s "+
			"(a bare hostname will be treated as a file path). Refusing to fall back to a local file",
		strings.Join(remoteSchemes, ", "),
	)
}
