package main

import (
	"context"
	"flag"
	"fmt"

	"soundboard-api/internal/db/gen"
)

// runApplyNames pushes the names file into an already-imported database.
//
// This is the revision path: correct a label in names.json, run this, and the grid picks
// it up without re-importing audio or disturbing play counts. Only clips whose stored
// name actually differs are written, so a no-op run reports nothing changed.
func runApplyNames(args []string) error {
	fs := flag.NewFlagSet("apply-names", flag.ExitOnError)
	namesPath := fs.String("names", defaultNamesPath, "path to the names file")
	dryRun := fs.Bool("dry-run", false, "show the renames without writing them")
	if err := fs.Parse(args); err != nil {
		return err
	}

	overrides, err := loadNameOverrides(*namesPath)
	if err != nil {
		return err
	}

	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	ctx := context.Background()

	soundbites, err := st.queries.ListSoundbites(ctx)
	if err != nil {
		return fmt.Errorf("list soundbites: %w", err)
	}

	var renamed, unchanged int
	for _, sb := range soundbites {
		want := nameFor(sb.Filename, overrides)
		if want == sb.Name {
			unchanged++
			continue
		}

		fmt.Printf("  %-42s %q -> %q\n", sb.Filename, sb.Name, want)
		renamed++

		if *dryRun {
			continue
		}
		if err := st.queries.UpdateSoundbiteName(ctx, gen.UpdateSoundbiteNameParams{
			Name:     want,
			Filename: sb.Filename,
		}); err != nil {
			return fmt.Errorf("rename %s: %w", sb.Filename, err)
		}
	}

	verb := "renamed"
	if *dryRun {
		verb = "would rename"
	}
	fmt.Printf("\n%s %d, left alone %d\n", verb, renamed, unchanged)
	return nil
}
