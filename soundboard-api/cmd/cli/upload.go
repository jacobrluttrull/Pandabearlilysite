package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"soundboard-api/internal/audio"
	"soundboard-api/internal/clipname"
	"soundboard-api/internal/db/gen"
)

// runUpload adds a single clip — the everyday path for "I made a new soundbite".
//
// Only the file is required. The clip's length is measured from the audio itself and its
// label is derived from the filename, so adding a clip is one command with one argument.
// Both can still be given explicitly when the defaults are not what you want.
func runUpload(args []string) error {
	fs := flag.NewFlagSet("upload", flag.ExitOnError)
	file := fs.String("file", "", "path to the audio file to add (required)")
	name := fs.String("name", "", "display name (default: derived from the filename)")
	dateMade := fs.String("date-made", "", "date the clip was made, e.g. 2026-08-27 (optional)")
	namesPath := fs.String("names", clipname.DefaultPath, "path to the names file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Allow `cli upload some-clip.mp3` as well as `cli upload -file some-clip.mp3`.
	if *file == "" && fs.NArg() == 1 {
		*file = fs.Arg(0)
	}
	if *file == "" {
		fs.Usage()
		return errors.New("-file is required")
	}

	if *dateMade != "" {
		if _, err := time.Parse("2006-01-02", *dateMade); err != nil {
			return fmt.Errorf("date-made must look like 2026-08-27: %w", err)
		}
	}

	seconds, err := audio.MP3Duration(*file)
	if err != nil {
		return fmt.Errorf("measure %s: %w", *file, err)
	}

	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	if err := os.MkdirAll(st.cfg.AudioDir, 0o755); err != nil {
		return fmt.Errorf("create audio dir: %w", err)
	}

	filename := filepath.Base(*file)
	destPath := filepath.Join(st.cfg.AudioDir, filename)

	overrides, err := clipname.LoadOverrides(*namesPath)
	if err != nil {
		return err
	}

	label := *name
	if label == "" {
		label = clipname.For(filename, overrides)
	}

	copied, err := copyFileIfMissing(*file, destPath)
	if err != nil {
		return fmt.Errorf("copy audio file: %w", err)
	}

	created, err := st.queries.CreateSoundbiteIfNew(context.Background(), gen.CreateSoundbiteIfNewParams{
		Name:          label,
		Filename:      filename,
		DateMade:      sql.NullString{String: *dateMade, Valid: *dateMade != ""},
		LengthSeconds: seconds,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s is already stored — use `cli rename` to relabel it", filename)
	}
	if err != nil {
		if copied {
			os.Remove(destPath)
		}
		return fmt.Errorf("save soundbite: %w", err)
	}

	// Record the label so it survives a later apply-names run.
	overrides[filename] = label
	if err := clipname.SaveOverrides(*namesPath, overrides); err != nil {
		return err
	}

	fmt.Printf("added %q as soundbite #%d (%s, %.1fs)\n", created.Name, created.ID, filename, seconds)
	return nil
}
