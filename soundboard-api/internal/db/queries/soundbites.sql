-- name: ListSoundbites :many
SELECT id, name, filename, date_stored, date_made, length_seconds, play_count
FROM soundbites
ORDER BY name COLLATE NOCASE ASC;

-- name: GetSoundbite :one
SELECT id, name, filename, date_stored, date_made, length_seconds, play_count
FROM soundbites
WHERE id = ?
LIMIT 1;

-- name: CreateSoundbite :one
INSERT INTO soundbites (name, filename, date_made, length_seconds)
VALUES (?, ?, ?, ?)
RETURNING id, name, filename, date_stored, date_made, length_seconds, play_count;

-- CreateSoundbiteIfNew inserts a clip unless one with the same filename is already
-- stored. On conflict no row is returned, so callers get sql.ErrNoRows and can treat
-- that as "already imported" rather than a failure. This is what makes a bulk import
-- safe to re-run after a partial failure.
-- name: CreateSoundbiteIfNew :one
INSERT INTO soundbites (name, filename, date_made, length_seconds)
VALUES (?, ?, ?, ?)
ON CONFLICT (filename) DO NOTHING
RETURNING id, name, filename, date_stored, date_made, length_seconds, play_count;

-- IncrementPlayCount bumps a clip's play tally by one and returns the new total.
--
-- The increment is computed inside SQLite rather than read-then-written in Go, so two
-- plays arriving at once cannot read the same starting value and lose a count. A
-- missing id yields no row, which the handler reports as a 404.
-- name: IncrementPlayCount :one
UPDATE soundbites
SET play_count = play_count + 1
WHERE id = ?
RETURNING play_count;

-- UpdateSoundbiteName re-labels a clip without touching its audio or play count.
-- Keyed by filename so the names file can be written against files on disk rather than
-- database ids the user never sees.
-- name: UpdateSoundbiteName :exec
UPDATE soundbites
SET name = ?
WHERE filename = ?;

-- DeleteSoundbite removes a clip's row. Keyed by filename to match the names file and
-- the audio dir; the caller is responsible for deleting the audio file itself.
-- name: DeleteSoundbite :execrows
DELETE FROM soundbites
WHERE filename = ?;

-- ResetAllPlayCounts zeroes every tally. Used to clear test plays before launch.
-- name: ResetAllPlayCounts :execrows
UPDATE soundbites
SET play_count = 0;

-- ResetPlayCount zeroes one clip's tally.
-- name: ResetPlayCount :execrows
UPDATE soundbites
SET play_count = 0
WHERE filename = ?;

-- SetSoundbiteDateMade records when a clip was originally made. Filenames carry no date,
-- so this is filled in by hand after import. Passing NULL clears it.
-- name: SetSoundbiteDateMade :execrows
UPDATE soundbites
SET date_made = ?
WHERE filename = ?;

-- SetPlayCount overwrites a clip's tally outright. Used when merging duplicates, so the
-- surviving clip keeps the plays its copies had rather than losing them.
-- name: SetPlayCount :execrows
UPDATE soundbites
SET play_count = ?
WHERE filename = ?;
