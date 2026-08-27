// Package seed brings the database in line with the clip audio shipped alongside it.
package seed

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"soundboard-api/internal/audio"
	"soundboard-api/internal/clipname"
	"soundboard-api/internal/db/gen"
)

// Clips inserts a row for every audio file that does not have one yet, and reports how
// many it added.
//
// This is what makes a deploy self-sufficient. Clip audio ships inside the container
// image while the database lives on a mounted volume, so a fresh volume starts empty and
// a new clip arrives with no row. Seeding on boot closes that gap: adding a soundbite is
// committing the file, with no migration, no admin endpoint, and no manual step against
// production.
//
// Existing rows are never modified. Play counts and any renames made through the CLI
// survive untouched, and a clip removed from the database stays removed rather than
// reappearing on the next restart — deleting its audio file is what makes removal stick.
func Clips(ctx context.Context, queries *gen.Queries, audioDir string, names map[string]string) (int, error) {
	filenames, err := listMP3s(audioDir)
	if err != nil {
		return 0, err
	}
	if len(filenames) == 0 {
		return 0, nil
	}

	stored, err := storedFilenames(ctx, queries)
	if err != nil {
		return 0, err
	}

	var added int
	for _, filename := range filenames {
		if stored[filename] {
			continue
		}

		seconds, err := audio.MP3Duration(filepath.Join(audioDir, filename))
		if err != nil {
			// One unreadable file must not stop the service from starting.
			log.Printf("seed: skipping %s: %v", filename, err)
			continue
		}

		name := clipname.For(filename, names)

		_, err = queries.CreateSoundbiteIfNew(ctx, gen.CreateSoundbiteIfNewParams{
			Name:          name,
			Filename:      filename,
			DateMade:      sql.NullString{},
			LengthSeconds: seconds,
		})
		if errors.Is(err, sql.ErrNoRows) {
			continue // already present
		}
		if err != nil {
			return added, fmt.Errorf("seed %s: %w", filename, err)
		}

		log.Printf("seed: added %q (%s, %.1fs)", name, filename, seconds)
		added++
	}

	return added, nil
}

func listMP3s(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		log.Printf("seed: audio dir %s does not exist, nothing to seed", dir)
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read audio dir: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ".mp3") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func storedFilenames(ctx context.Context, queries *gen.Queries) (map[string]bool, error) {
	soundbites, err := queries.ListSoundbites(ctx)
	if err != nil {
		return nil, fmt.Errorf("list soundbites: %w", err)
	}

	stored := make(map[string]bool, len(soundbites))
	for _, sb := range soundbites {
		stored[sb.Filename] = true
	}
	return stored, nil
}
