package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// findMP3s returns the names of the .mp3 files directly inside dir, sorted so an import
// run reports clips in a predictable order.
//
// Subdirectories are deliberately not searched: soundbites.filename is UNIQUE, so two
// clips sharing a basename in different folders would collide, and the second would be
// silently treated as a duplicate of the first. Flattening that decision is the
// importer's job to avoid, not to guess at.
func findMP3s(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read clip folder: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".mp3") {
			continue
		}
		names = append(names, entry.Name())
	}

	sort.Strings(names)
	return names, nil
}
