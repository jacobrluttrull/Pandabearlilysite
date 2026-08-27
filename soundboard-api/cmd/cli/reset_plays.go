package main

import (
	"context"
	"flag"
	"fmt"
)

// runResetPlays zeroes play tallies — all of them, or just the clips named.
//
// The main use is clearing counts accumulated while testing, so the site launches from
// zero rather than showing numbers nobody earned.
func runResetPlays(args []string) error {
	fs := flag.NewFlagSet("reset-plays", flag.ExitOnError)
	all := fs.Bool("all", false, "reset every clip (required if no filenames are given)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	filenames := fs.Args()
	if len(filenames) == 0 && !*all {
		fs.Usage()
		return fmt.Errorf("pass filenames to reset, or -all to reset everything")
	}
	if len(filenames) > 0 && *all {
		return fmt.Errorf("pass either filenames or -all, not both")
	}

	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	ctx := context.Background()

	if *all {
		rows, err := st.queries.ResetAllPlayCounts(ctx)
		if err != nil {
			return fmt.Errorf("reset play counts: %w", err)
		}
		fmt.Printf("reset play counts on %d clip(s)\n", rows)
		return nil
	}

	var reset, missing int
	for _, filename := range filenames {
		rows, err := st.queries.ResetPlayCount(ctx, filename)
		if err != nil {
			return fmt.Errorf("reset %s: %w", filename, err)
		}
		if rows == 0 {
			fmt.Printf("  not stored: %s\n", filename)
			missing++
			continue
		}
		fmt.Printf("  reset %s\n", filename)
		reset++
	}

	fmt.Printf("\nreset %d, not stored %d\n", reset, missing)
	return nil
}
