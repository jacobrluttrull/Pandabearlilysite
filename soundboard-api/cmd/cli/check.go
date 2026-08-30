package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sort"

	"soundboard-api/internal/clipname"
)

// runCheck audits the places a clip exists — the database, the clip store, and the names
// file — and reports anything that has drifted out of sync.
//
// The clip store is the authority on what a visitor can actually hear, which since the
// move to R2 is the bucket rather than a folder on this machine. Auditing the local
// folder instead would report every clip published from another checkout as missing
// audio, when the site serves it perfectly well.
//
// The local folder is still reported, but as a working copy: it is the second copy of the
// audio and worth knowing about, not a fault when it lags.
func runCheck(args []string) error {
	fs_ := flag.NewFlagSet("check", flag.ExitOnError)
	namesPath := fs_.String("names", clipname.DefaultPath, "path to the names file")
	if err := fs_.Parse(args); err != nil {
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

	overrides, err := clipname.LoadOverrides(*namesPath)
	if err != nil {
		return err
	}

	// What the site can actually serve.
	publishedList, err := st.clips.List(ctx)
	if err != nil {
		return fmt.Errorf("list %s: %w", st.clips.Describe(), err)
	}

	// The local working copy. Absent is fine — a checkout that publishes straight to the
	// bucket never needs one.
	onDisk, err := localClips(st.cfg.AudioDir)
	if err != nil {
		return err
	}

	inDB := make(map[string]bool, len(soundbites))
	for _, sb := range soundbites {
		inDB[sb.Filename] = true
	}
	published := make(map[string]bool, len(publishedList))
	for _, filename := range publishedList {
		published[filename] = true
	}
	local := make(map[string]bool, len(onDisk))
	for _, filename := range onDisk {
		local[filename] = true
	}

	problems := 0

	// Rows whose audio is missing from the store: the clip is listed on the site but
	// 404s when tapped. This is the one that matters.
	for _, sb := range soundbites {
		if !published[sb.Filename] {
			fmt.Printf("  MISSING AUDIO  %s (row #%d) — listed but not in %s\n",
				sb.Filename, sb.ID, st.clips.Describe())
			problems++
		}
	}

	// Audio with no row: dead weight, never served, still paying for storage.
	for _, filename := range publishedList {
		if !inDB[filename] {
			fmt.Printf("  ORPHAN CLIP    %s — in %s but not in the database\n",
				filename, st.clips.Describe())
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

	// Byte-identical clips that dedupe would collapse. This needs the bytes, so it can
	// only cover what is on this machine; clips published from elsewhere are skipped and
	// counted, so a partial scan is never mistaken for a clean one.
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

	fmt.Printf("\n%d clip(s) in the database, %d in %s, %d name(s) on file\n",
		len(soundbites), len(publishedList), st.clips.Describe(), len(overrides))

	// The local folder is a backup, not a fault. Say where it stands without counting it
	// as a problem, since losing both it and the store is what actually loses the clips.
	missingLocally := 0
	for _, sb := range soundbites {
		if published[sb.Filename] && !local[sb.Filename] {
			missingLocally++
		}
	}
	if missingLocally > 0 {
		fmt.Printf("%d published clip(s) have no copy in %s "+
			"(not a problem, but the store is then the only copy of those)\n",
			missingLocally, st.cfg.AudioDir)
	}

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

// localClips lists the local working copy, treating an absent folder as empty rather than
// an error: publishing straight to a bucket never creates one.
func localClips(dir string) ([]string, error) {
	files, err := findMP3s(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return files, nil
}
