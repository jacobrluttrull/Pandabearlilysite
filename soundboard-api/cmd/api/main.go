package main

import (
	"log"
	"net/http"
	"os"

	"soundboard-api/internal/api"
	"soundboard-api/internal/config"
	"soundboard-api/internal/db"
	"soundboard-api/internal/db/gen"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.AudioDir, 0o755); err != nil {
		log.Fatalf("create audio dir: %v", err)
	}

	sqlDB, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer sqlDB.Close()

	if err := db.Migrate(sqlDB); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	handler := api.NewHandler(gen.New(sqlDB), cfg.AudioDir, cfg.AllowedOrigin)

	log.Printf("soundboard-api listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
