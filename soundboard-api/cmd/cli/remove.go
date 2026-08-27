package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"soundboard-api/internal/clipname"
)

// runRemove deletes clips by filename: the database row, the stored audio file, and the
// entry in the names file, so a removed clip leaves nothing behind to re-import.
//
// The source file in the staging folder is left alone — removing a clip from the site
// should not delete the original off your disk.
func runRemove(args []string) error {
	fs := flag.NewFlagSet("remove", flag.ExitOnError)
	namesPath := fs.String("names", clipname.DefaultPath, "path to the names file")
	dryRun := fs.Bool("dry-run", false, "show what would be removed without deleting anything")
	if err := fs.Parse(args); err != nil {
		return err
	}

	filenames := fs.Args()
	if len(filenames) == 0 {
		fs.Usage()
		return errors.New("pass one or more filenames to remove, e.g. cli remove \"thats-enough-head (1).mp3\"")
	}

	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	overrides, err := clipname.LoadOverrides(*namesPath)
	if err != nil {
		return err
	}

	ctx := context.Background()

	var removed, missing int
	for _, filename := range filenames {
		if *dryRun {
			fmt.Printf("  would remove %s\n", filename)
			removed++
			continue
		}

		rows, err := st.queries.DeleteSoundbite(ctx, filename)
		if err != nil {
			return fmt.Errorf("delete %s: %w", filename, err)
		}
		if rows == 0 {
			fmt.Fprintf(os.Stderr, "  not stored: %s\n", filename)
			missing++
			continue
		}

		audioPath := filepath.Join(st.cfg.AudioDir, filename)
		if err := os.Remove(audioPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("delete audio %s: %w", audioPath, err)
		}

		delete(overrides, filename)
		removed++
		fmt.Printf("  removed %s\n", filename)
	}

	if !*dryRun && removed > 0 {
		if err := clipname.SaveOverrides(*namesPath, overrides); err != nil {
			return err
		}
	}

	fmt.Printf("\nremoved %d, not stored %d\n", removed, missing)
	return nil
}
