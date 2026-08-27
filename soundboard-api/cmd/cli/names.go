package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// runNames scaffolds or refreshes the names file, seeding an entry for every clip in a
// folder using the name derived from its filename.
//
// Existing entries are never overwritten, so this is safe to re-run after adding clips:
// new files get a starting point, and names already corrected by hand are left alone.
func runNames(args []string) error {
	fs := flag.NewFlagSet("names", flag.ExitOnError)
	dir := fs.String("dir", "", "folder of .mp3 clips to seed names for (required)")
	namesPath := fs.String("names", defaultNamesPath, "path to the names file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *dir == "" {
		fs.Usage()
		return errors.New("-dir is required")
	}

	filenames, err := findMP3s(*dir)
	if err != nil {
		return err
	}
	if len(filenames) == 0 {
		return fmt.Errorf("no .mp3 files found in %s", *dir)
	}

	overrides, err := loadNameOverrides(*namesPath)
	if err != nil {
		return err
	}

	var added, kept int
	for _, filename := range filenames {
		if _, exists := overrides[filename]; exists {
			kept++
			continue
		}
		overrides[filename] = displayName(filename)
		added++
	}

	// Entries whose file is gone are worth surfacing but not deleting: a clip may just
	// be staged elsewhere, and silently dropping a curated name is worse than a warning.
	present := make(map[string]bool, len(filenames))
	for _, filename := range filenames {
		present[filename] = true
	}
	for filename := range overrides {
		if !present[filename] {
			fmt.Fprintf(os.Stderr, "  note: %s is named but not in %s\n", filename, *dir)
		}
	}

	if err := saveNameOverrides(*namesPath, overrides); err != nil {
		return err
	}

	fmt.Printf("wrote %s: %d new entr(ies), %d left as written\n", *namesPath, added, kept)
	fmt.Println("edit the values, then run `cli import` (or `cli apply-names` if already imported)")
	return nil
}
