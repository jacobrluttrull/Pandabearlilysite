package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"soundboard-api/internal/clipname"
	"soundboard-api/internal/db/gen"
)

// runDedupe finds clips whose audio is byte-for-byte identical and collapses each group
// down to one, folding the copies' play counts into the survivor.
//
// Duplicates come from file managers: copying a clip yields "clip (1).mp3", which is a
// different filename and so passes the UNIQUE constraint and imports as a separate clip.
// Content hashing is what actually catches them.
func runDedupe(args []string) error {
	fs := flag.NewFlagSet("dedupe", flag.ExitOnError)
	namesPath := fs.String("names", clipname.DefaultPath, "path to the names file")
	dryRun := fs.Bool("dry-run", false, "show the duplicate groups without removing anything")
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

	groups, err := groupByContent(soundbites, st.cfg.AudioDir)
	if err != nil {
		return err
	}

	if len(groups) == 0 {
		fmt.Println("no duplicate clips found")
		return nil
	}

	overrides, err := clipname.LoadOverrides(*namesPath)
	if err != nil {
		return err
	}

	var removed int
	for _, group := range groups {
		keeper, copies := pickKeeper(group)

		mergedPlays := keeper.PlayCount
		for _, dup := range copies {
			mergedPlays += dup.PlayCount
		}

		fmt.Printf("  keep   %s (%d plays after merge)\n", keeper.Filename, mergedPlays)
		for _, dup := range copies {
			fmt.Printf("  remove %s (%d plays)\n", dup.Filename, dup.PlayCount)
		}

		if *dryRun {
			removed += len(copies)
			continue
		}

		if mergedPlays != keeper.PlayCount {
			if _, err := st.queries.SetPlayCount(ctx, gen.SetPlayCountParams{
				PlayCount: mergedPlays,
				Filename:  keeper.Filename,
			}); err != nil {
				return fmt.Errorf("merge plays into %s: %w", keeper.Filename, err)
			}
		}

		for _, dup := range copies {
			if _, err := st.queries.DeleteSoundbite(ctx, dup.Filename); err != nil {
				return fmt.Errorf("delete %s: %w", dup.Filename, err)
			}

			audioPath := filepath.Join(st.cfg.AudioDir, dup.Filename)
			if err := os.Remove(audioPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("delete audio %s: %w", audioPath, err)
			}

			delete(overrides, dup.Filename)
			removed++
		}
	}

	if !*dryRun && removed > 0 {
		if err := clipname.SaveOverrides(*namesPath, overrides); err != nil {
			return err
		}
	}

	verb := "removed"
	if *dryRun {
		verb = "would remove"
	}
	fmt.Printf("\n%s %d duplicate(s) across %d group(s)\n", verb, removed, len(groups))

	if !*dryRun && removed > 0 {
		fmt.Println("note: the source files are still in your clips folder — delete them there")
		fmt.Println("      too, or the next import will add them back")
	}
	return nil
}

// groupByContent returns only those sets of clips that share identical audio. Clips whose
// audio file is missing are skipped rather than treated as matching each other.
func groupByContent(soundbites []gen.Soundbite, audioDir string) ([][]gen.Soundbite, error) {
	byHash := map[string][]gen.Soundbite{}

	for _, sb := range soundbites {
		path := filepath.Join(audioDir, sb.Filename)

		hash, err := fileHash(path)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "  note: %s has no audio file, skipping\n", sb.Filename)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("hash %s: %w", sb.Filename, err)
		}
		byHash[hash] = append(byHash[hash], sb)
	}

	var groups [][]gen.Soundbite
	for _, group := range byHash {
		if len(group) > 1 {
			groups = append(groups, group)
		}
	}

	// Stable output so repeated runs and dry-run previews agree.
	sort.Slice(groups, func(i, j int) bool {
		return groups[i][0].Filename < groups[j][0].Filename
	})
	return groups, nil
}

// pickKeeper chooses which clip in a duplicate group survives: the shortest filename,
// breaking ties alphabetically. That reliably prefers "clip.mp3" over the copy-suffixed
// "clip (1).mp3" without having to pattern-match whatever suffix the file manager used.
func pickKeeper(group []gen.Soundbite) (keeper gen.Soundbite, copies []gen.Soundbite) {
	sorted := make([]gen.Soundbite, len(group))
	copy(sorted, group)

	sort.Slice(sorted, func(i, j int) bool {
		if len(sorted[i].Filename) != len(sorted[j].Filename) {
			return len(sorted[i].Filename) < len(sorted[j].Filename)
		}
		return sorted[i].Filename < sorted[j].Filename
	})

	return sorted[0], sorted[1:]
}
