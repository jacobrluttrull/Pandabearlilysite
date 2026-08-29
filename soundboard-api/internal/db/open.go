package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

// remoteMaxOpenConns caps concurrent connections to a remote database. It is a ceiling
// against a runaway request spike opening sockets without limit, not a target — normal
// fan-site traffic will sit far below it.
const remoteMaxOpenConns = 10

// Open connects to the soundboard database described by dsn.
//
// dsn is either a filesystem path, which opens (and creates) a local SQLite file for
// development, or a libsql:// URL, which opens a remote Turso database for production.
// Callers pass the configured value straight through; the shape of the string decides
// which of the two it is.
func Open(dsn string) (*sql.DB, error) {
	driverName, connDSN := driverFor(dsn)

	db, err := sql.Open(driverName, connDSN)
	if err != nil {
		return nil, fmt.Errorf("open %s db: %w", driverName, err)
	}

	if IsRemoteDSN(dsn) {
		// The single-connection rule below is about a file on disk, and does not carry
		// over here: Turso is a network service that serialises writes on its own side,
		// so a one-connection pool would buy no extra safety and would instead queue
		// every request behind one in-flight HTTP round trip. Latency, not lock
		// contention, is the limit on a remote database — so let requests overlap.
		db.SetMaxOpenConns(remoteMaxOpenConns)
		db.SetMaxIdleConns(remoteMaxOpenConns)

		// Idle connections to a remote host go stale silently — a load balancer or the
		// server can drop one without us noticing until a query fails on it. Retiring
		// them on a timer costs a reconnect on a quiet site and avoids that.
		db.SetConnMaxIdleTime(5 * time.Minute)
	} else {
		// SQLite only supports one writer at a time; a single connection avoids
		// SQLITE_BUSY errors under concurrent access. This also serialises reads, which
		// is fine at fan-site traffic — if plays ever outgrow it, the next step is
		// separate read and write pools rather than a bigger single pool.
		db.SetMaxOpenConns(1)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping %s db: %w", driverName, err)
	}

	return db, nil
}
