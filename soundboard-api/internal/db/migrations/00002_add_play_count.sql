-- +goose Up
-- Tracks how many times a clip has been played. Existing rows start at zero rather
-- than NULL so the API can always return a number and callers never special-case it.
ALTER TABLE soundbites ADD COLUMN play_count INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE soundbites DROP COLUMN play_count;
