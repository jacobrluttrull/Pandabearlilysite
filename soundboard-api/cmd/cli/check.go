package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
)

// runCheck audits the three places a clip exists — the database, the audio dir, and the
// names file — and reports anything that has drifted out of sync.
//
// Worth running before launch: it catches rows whose audio never copied, audio files
// orphaned by a manual delete, and duplicate content that dedupe would collapse.
func runCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	namesPath := fs.String("names", defaultNamesPath, "path to the names file")
	if err := fs.Parse(args); err != nil {
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

	overrides, err := loadNameOverrides(*namesPath)
	if err != nil {
		return err
	}

	onDisk, err := findMP3s(st.cfg.AudioDir)
	if err != nil {
		return err
	}

	inDB := make(map[string]bool, len(soundbites))
	for _, sb := range soundbites {
		inDB[sb.Filename] = true
	}
	stored := make(map[string]bool, len(onDisk))
	for _, filename := range onDisk {
		stored[filename] = true
	}

	problems := 0

	// Rows whose audio is missing: the clip is listed but cannot play.
	for _, sb := range soundbites {
		if !stored[sb.Filename] {
			fmt.Printf("  MISSING AUDIO  %s (row #%d) — listed but no file in %s\n",
				sb.Filename, sb.ID, st.cfg.AudioDir)
			problems++
		}
	}

	// Audio with no row: dead weight, never served.
	for _, filename := range onDisk {
		if !inDB[filename] {
			fmt.Printf("  ORPHAN FILE    %s — in %s but not in the database\n",
				filename, st.cfg.AudioDir)
			problems++
		}
	}

	// Named but absent: usually a leftover after removing a clip.
	var strayNames []string
	for filename := range overrides {
		if !inDB[filename] {
			strayNames = append(strayNames, filename)
		}
	}
	sort.Strings(strayNames)
	for _, filename := range strayNames {
		fmt.Printf("  STRAY NAME     %s — named in %s but not stored\n", filename, *namesPath)
		problems++
	}

	// Byte-identical clips that dedupe would collapse.
	groups, err := groupByContent(soundbites, st.cfg.AudioDir)
	if err != nil {
		return err
	}
	for _, group := range groups {
		names := make([]string, len(group))
		for i, sb := range group {
			names[i] = sb.Filename
		}
		fmt.Printf("  DUPLICATE      identical audio: %v — run `cli dedupe`\n", names)
		problems += len(group) - 1
	}

	fmt.Printf("\n%d clip(s) in the database, %d audio file(s) on disk, %d name(s) on file\n",
		len(soundbites), len(onDisk), len(overrides))

	// Not an error condition, just worth knowing before launch.
	var undated int
	for _, sb := range soundbites {
		if !sb.DateMade.Valid {
			undated++
		}
	}
	if undated > 0 {
		fmt.Printf("%d clip(s) have no date_made (optional — set with `cli set-date`)\n", undated)
	}

	if problems == 0 {
		fmt.Println("no problems found")
		return nil
	}
	fmt.Fprintf(os.Stderr, "%d problem(s) found\n", problems)
	return nil
}
