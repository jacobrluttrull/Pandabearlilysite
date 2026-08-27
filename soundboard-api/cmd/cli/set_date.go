package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"time"

	"soundboard-api/internal/db/gen"
)

// runSetDate records when a clip was originally made.
//
// Filenames carry no date, so import leaves date_made empty and it is filled in here.
// Passing -clear removes the date again.
func runSetDate(args []string) error {
	fs := flag.NewFlagSet("set-date", flag.ExitOnError)
	clear := fs.Bool("clear", false, "clear the date instead of setting one")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var filename, date string
	switch {
	case *clear && fs.NArg() == 1:
		filename = fs.Arg(0)
	case !*clear && fs.NArg() == 2:
		filename, date = fs.Arg(0), fs.Arg(1)
	default:
		fs.Usage()
		return errors.New("usage: cli set-date <filename.mp3> <YYYY-MM-DD>   (or -clear <filename.mp3>)")
	}

	if !*clear {
		// Validated up front so a typo lands as an error here rather than as a date the
		// frontend cannot parse.
		if _, err := time.Parse("2006-01-02", date); err != nil {
			return fmt.Errorf("date must look like 2026-08-27: %w", err)
		}
	}

	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	value := sql.NullString{String: date, Valid: !*clear}

	rows, err := st.queries.SetSoundbiteDateMade(context.Background(), gen.SetSoundbiteDateMadeParams{
		DateMade: value,
		Filename: filename,
	})
	if err != nil {
		return fmt.Errorf("set date on %s: %w", filename, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s is not stored — run `cli list` to see what is", filename)
	}

	if *clear {
		fmt.Printf("cleared date on %s\n", filename)
		return nil
	}
	fmt.Printf("%s made %s\n", filename, date)
	return nil
}
