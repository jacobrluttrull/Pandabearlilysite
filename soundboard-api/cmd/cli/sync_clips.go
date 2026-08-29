package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// runSyncClips uploads local clip files that the store does not have yet.
//
// This exists for two moments. The first is the move to R2: every clip already had a
// database row and a file on disk, and only the bytes needed to reach the bucket. The
// second is repair — if an upload failed halfway, or a clip was restored from a backup,
// this puts the store back in step without touching the database at all.
//
// It only ever adds. A clip present in the store but missing locally is left alone,
// because the store is the published copy and this machine's folder is not authoritative
// over it.
func runSyncClips(args []string) error {
	fs := flag.NewFlagSet("sync-clips", flag.ExitOnError)
	dir := fs.String("dir", "", "folder of .mp3 files to publish (default: the configured audio dir)")
	dryRun := fs.Bool("dry-run", false, "list what would be uploaded without uploading")
	if err := fs.Parse(args); err != nil {
		return err
	}

	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	srcDir := *dir
	if srcDir == "" {
		srcDir = st.cfg.AudioDir
	}

	local, err := findMP3s(srcDir)
	if err != nil {
		return err
	}
	if len(local) == 0 {
		return fmt.Errorf("no .mp3 files found in %s", srcDir)
	}

	ctx := context.Background()

	remote, err := st.clips.List(ctx)
	if err != nil {
		return err
	}
	have := make(map[string]bool, len(remote))
	for _, name := range remote {
		have[name] = true
	}

	fmt.Printf("%s\n  local:  %d file(s) in %s\n  stored: %d\n\n",
		st.clips.Describe(), len(local), srcDir, len(remote))

	var uploaded, skipped int
	for _, filename := range local {
		if have[filename] {
			skipped++
			continue
		}
		if *dryRun {
			fmt.Printf("  would upload  %s\n", filename)
			uploaded++
			continue
		}
		if err := publishClip(ctx, st, filepath.Join(srcDir, filename), filename); err != nil {
			fmt.Fprintf(os.Stderr, "  FAIL  %-40s %v\n", filename, err)
			continue
		}
		fmt.Printf("  uploaded  %s\n", filename)
		uploaded++
	}

	verb := "uploaded"
	if *dryRun {
		verb = "would upload"
	}
	fmt.Printf("\n%s %d, already stored %d\n", verb, uploaded, skipped)
	return nil
}
