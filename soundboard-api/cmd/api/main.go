package main

import (
	"context"
	"log"
	"net/http"

	"soundboard-api/internal/api"
	"soundboard-api/internal/clipname"
	"soundboard-api/internal/config"
	"soundboard-api/internal/db"
	"soundboard-api/internal/db/gen"
	"soundboard-api/internal/seed"
)

func main() {
	cfg := config.Load()

	sqlDB, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer sqlDB.Close()

	if err := db.Migrate(sqlDB); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	queries := gen.New(sqlDB)

	// Clip audio ships in the image while the database lives on a volume, so a fresh
	// volume — or a newly added clip — needs rows creating before anything can be
	// served. Existing rows, including play counts, are left alone.
	names, err := clipname.LoadOverrides(clipname.DefaultPath)
	if err != nil {
		log.Printf("seed: continuing without name overrides: %v", err)
		names = map[string]string{}
	}

	added, err := seed.Clips(context.Background(), queries, cfg.AudioDir, names)
	if err != nil {
		log.Fatalf("seed clips: %v", err)
	}
	if added > 0 {
		log.Printf("seed: added %d clip(s) from %s", added, cfg.AudioDir)
	}

	handler := api.NewHandler(queries, api.Options{
		AudioDir:      cfg.AudioDir,
		AllowedOrigin: cfg.AllowedOrigin,
		StaticDir:     cfg.StaticDir,
	})

	if cfg.StaticDir != "" {
		log.Printf("serving frontend from %s", cfg.StaticDir)
	}
	log.Printf("soundboard-api listening on %s", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
