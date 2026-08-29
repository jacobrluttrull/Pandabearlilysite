package main

import (
	"log"
	"net/http"

	"soundboard-api/internal/api"
	"soundboard-api/internal/clipstore"
	"soundboard-api/internal/config"
	"soundboard-api/internal/db"
	"soundboard-api/internal/db/gen"
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

	// Clip audio lives in R2 in production and in a local folder during development.
	// There is no seeding step any more: the CLI writes the audio and its database row
	// together, so a row without its audio cannot arrive here in the first place.
	if err := clipstore.RequireComplete(cfg.R2()); err != nil {
		log.Fatalf("%v", err)
	}

	clips := clipstore.Open(cfg.R2(), cfg.AudioDir)
	log.Printf("clips: %s", clips.Describe())

	handler := api.NewHandler(queries, api.Options{
		Clips:         clips,
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
