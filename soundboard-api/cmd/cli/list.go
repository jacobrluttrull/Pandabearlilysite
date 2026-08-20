package main

import (
	"context"
	"fmt"

	"soundboard-api/internal/config"
	sbdb "soundboard-api/internal/db"
	"soundboard-api/internal/db/gen"
)

func runList(args []string) error {
	cfg := config.Load()

	sqlDB, err := sbdb.Open(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer sqlDB.Close()

	if err := sbdb.Migrate(sqlDB); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	queries := gen.New(sqlDB)
	soundbites, err := queries.ListSoundbites(context.Background())
	if err != nil {
		return fmt.Errorf("list soundbites: %w", err)
	}

	if len(soundbites) == 0 {
		fmt.Println("no soundbites stored yet")
		return nil
	}

	for _, sb := range soundbites {
		dateMade := "-"
		if sb.DateMade.Valid {
			dateMade = sb.DateMade.String
		}
		fmt.Printf("#%-4d %-30s %6.1fs  made:%-12s file:%s\n", sb.ID, sb.Name, sb.LengthSeconds, dateMade, sb.Filename)
	}
	return nil
}
