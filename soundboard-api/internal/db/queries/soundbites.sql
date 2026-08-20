-- name: ListSoundbites :many
SELECT id, name, filename, date_stored, date_made, length_seconds
FROM soundbites
ORDER BY name COLLATE NOCASE ASC;

-- name: GetSoundbite :one
SELECT id, name, filename, date_stored, date_made, length_seconds
FROM soundbites
WHERE id = ?
LIMIT 1;

-- name: CreateSoundbite :one
INSERT INTO soundbites (name, filename, date_made, length_seconds)
VALUES (?, ?, ?, ?)
RETURNING id, name, filename, date_stored, date_made, length_seconds;
