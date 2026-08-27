package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// fileHash returns the SHA-256 of a file's contents.
//
// Duplicate clips are detected by content rather than by name or duration: files copied
// by a file manager pick up names like "clip (1).mp3" that no name rule would connect,
// while two genuinely different clips can easily share a duration.
func fileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}
