-- +goose Up
CREATE TABLE soundbites (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    name           TEXT NOT NULL,
    filename       TEXT NOT NULL UNIQUE,
    date_stored    TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    date_made      TEXT,
    length_seconds REAL NOT NULL
);

-- +goose Down
DROP TABLE soundbites;
