package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"soundboard-api/internal/config"
	sbdb "soundboard-api/internal/db"
	"soundboard-api/internal/db/gen"
)

func runUpload(args []string) error {
	fs := flag.NewFlagSet("upload", flag.ExitOnError)
	name := fs.String("name", "", "display name shown on the soundboard (required)")
	file := fs.String("file", "", "path to the audio file to upload (required)")
	dateMade := fs.String("date-made", "", "date the clip was made, e.g. 2026-08-19 (optional)")
	length := fs.Float64("length", 0, "clip length in seconds (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" || *file == "" || *length <= 0 {
		fs.Usage()
		return fmt.Errorf("-name, -file, and -length are required")
	}

	cfg := config.Load()

	if err := os.MkdirAll(cfg.AudioDir, 0o755); err != nil {
		return fmt.Errorf("create audio dir: %w", err)
	}

	sqlDB, err := sbdb.Open(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer sqlDB.Close()

	if err := sbdb.Migrate(sqlDB); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	filename := filepath.Base(*file)
	destPath := filepath.Join(cfg.AudioDir, filename)

	if err := copyFile(*file, destPath); err != nil {
		return fmt.Errorf("copy audio file: %w", err)
	}

	queries := gen.New(sqlDB)
	created, err := queries.CreateSoundbite(context.Background(), gen.CreateSoundbiteParams{
		Name:          *name,
		Filename:      filename,
		DateMade:      sql.NullString{String: *dateMade, Valid: *dateMade != ""},
		LengthSeconds: *length,
	})
	if err != nil {
		os.Remove(destPath)
		return fmt.Errorf("save soundbite: %w", err)
	}

	fmt.Printf("uploaded %q as soundbite #%d (%s)\n", created.Name, created.ID, filename)
	return nil
}
