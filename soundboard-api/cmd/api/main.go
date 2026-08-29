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

	// Logged before the connection is attempted, so a misconfigured deploy says which
	// database it fell back to even when opening it then fails.
	log.Printf("database: %s", cfg.DatabaseDescription())

	// TURSO_DATABASE_URL being set states an intent to use a remote database. If the
	// value is malformed it would otherwise be taken as a filesystem path and the
	// service would start quietly on a local file that no deploy survives.
	if cfg.DatabaseURL != "" {
		if err := db.RequireRemote(cfg.DatabaseDSN()); err != nil {
			log.Fatalf("%v", err)
		}
	}

	sqlDB, err := db.Open(cfg.DatabaseDSN())
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
		AuthUser:      cfg.AuthUser,
		AuthPassword:  cfg.AuthPassword,
	})

	// Stated at boot for the same reason the database is: the unprotected mode is the
	// silent one, and this site is not meant to be reachable without a password.
	if cfg.AuthPassword == "" {
		log.Printf("auth: DISABLED — every visitor gets in. Set SOUNDBOARD_AUTH_PASSWORD before deploying.")
	} else {
		log.Printf("auth: password required (user %q)", cfg.AuthUser)
	}

	if cfg.StaticDir != "" {
		log.Printf("serving frontend from %s", cfg.StaticDir)
	}
	log.Printf("soundboard-api listening on %s", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
