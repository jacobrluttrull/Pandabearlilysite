package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"soundboard-api/internal/db/gen"
)

// runRename relabels one clip without opening the names file by hand.
//
// The new label is written to both the database and names.json, so the two never drift:
// a rename made here survives a later `apply-names`, and a rename made by editing the
// file survives a later import.
func runRename(args []string) error {
	fs := flag.NewFlagSet("rename", flag.ExitOnError)
	namesPath := fs.String("names", defaultNamesPath, "path to the names file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 2 {
		fs.Usage()
		return errors.New(`usage: cli rename <filename.mp3> "new display name"`)
	}
	filename, newName := fs.Arg(0), fs.Arg(1)

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

	var current string
	found := false
	for _, sb := range soundbites {
		if sb.Filename == filename {
			current, found = sb.Name, true
			break
		}
	}
	if !found {
		return fmt.Errorf("%s is not stored — run `cli list` to see what is", filename)
	}

	if err := st.queries.UpdateSoundbiteName(ctx, gen.UpdateSoundbiteNameParams{
		Name:     newName,
		Filename: filename,
	}); err != nil {
		return fmt.Errorf("rename %s: %w", filename, err)
	}

	overrides, err := loadNameOverrides(*namesPath)
	if err != nil {
		return err
	}
	overrides[filename] = newName
	if err := saveNameOverrides(*namesPath, overrides); err != nil {
		return err
	}

	fmt.Printf("%s: %q -> %q\n", filename, current, newName)
	fmt.Printf("also recorded in %s\n", *namesPath)
	return nil
}
