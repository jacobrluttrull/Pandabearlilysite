package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"soundboard-api/internal/audio"
	"soundboard-api/internal/config"
	sbdb "soundboard-api/internal/db"
	"soundboard-api/internal/db/gen"
)

// runImport bulk-loads a folder of .mp3 clips: it measures each clip's length, derives
// a display name from its filename, copies the audio into the configured audio dir, and
// inserts a row per clip in a single transaction.
//
// The command is idempotent. Clips already stored are skipped rather than failing, so a
// run interrupted halfway can simply be run again.
func runImport(args []string) error {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	dir := fs.String("dir", "", "folder of .mp3 clips to import (required)")
	dryRun := fs.Bool("dry-run", false, "report what would be imported without copying files or writing to the database")
	namesPath := fs.String("names", defaultNamesPath, "path to the names file used to override derived labels")
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

	cfg := config.Load()

	sqlDB, err := sbdb.Open(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer sqlDB.Close()

	if err := sbdb.Migrate(sqlDB); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	ctx := context.Background()

	stored, err := storedFilenames(ctx, gen.New(sqlDB))
	if err != nil {
		return err
	}

	if *dryRun {
		return reportDryRun(*dir, filenames, stored, overrides)
	}

	if err := os.MkdirAll(cfg.AudioDir, 0o755); err != nil {
		return fmt.Errorf("create audio dir: %w", err)
	}

	fmt.Printf("importing %d mp3 file(s) from %s\n\n", len(filenames), *dir)

	result, err := importClips(ctx, sqlDB, cfg, *dir, filenames, stored, overrides)
	if err != nil {
		return err
	}

	fmt.Printf("\nadded %d, skipped %d (already stored), failed %d\n",
		result.added, result.skipped, len(result.failures))
	for _, failure := range result.failures {
		fmt.Fprintf(os.Stderr, "  failed: %s: %v\n", failure.filename, failure.err)
	}
	return nil
}

type importResult struct {
	added    int
	skipped  int
	failures []importFailure
}

type importFailure struct {
	filename string
	err      error
}

// importClips does the work of a real (non-dry-run) import inside one transaction.
//
// Audio files are copied as clips are processed, but a clip that fails to measure or
// insert does not abort the batch — the failure is collected and the run continues, so
// one bad file out of hundreds does not cost you the whole import. If the transaction
// itself cannot commit, files copied during this run are removed again so the audio dir
// does not drift out of sync with the database.
func importClips(
	ctx context.Context,
	sqlDB *sql.DB,
	cfg config.Config,
	srcDir string,
	filenames []string,
	stored map[string]bool,
	overrides map[string]string,
) (importResult, error) {
	var result importResult

	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return result, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	queries := gen.New(sqlDB).WithTx(tx)

	var copiedPaths []string
	cleanupCopies := func() {
		for _, path := range copiedPaths {
			os.Remove(path)
		}
	}

	for _, filename := range filenames {
		if stored[filename] {
			result.skipped++
			fmt.Printf("  skip  %-40s already stored\n", filename)
			continue
		}

		srcPath := filepath.Join(srcDir, filename)

		seconds, err := audio.MP3Duration(srcPath)
		if err != nil {
			result.failures = append(result.failures, importFailure{filename, err})
			fmt.Printf("  FAIL  %-40s %v\n", filename, err)
			continue
		}

		destPath := filepath.Join(cfg.AudioDir, filename)
		copied, err := copyFileIfMissing(srcPath, destPath)
		if err != nil {
			result.failures = append(result.failures, importFailure{filename, err})
			fmt.Printf("  FAIL  %-40s %v\n", filename, err)
			continue
		}
		if copied {
			copiedPaths = append(copiedPaths, destPath)
		}

		name := nameFor(filename, overrides)
		_, err = queries.CreateSoundbiteIfNew(ctx, gen.CreateSoundbiteIfNewParams{
			Name:          name,
			Filename:      filename,
			DateMade:      sql.NullString{}, // not derivable from a filename; fill in later
			LengthSeconds: seconds,
		})
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// Another row already claimed this filename; treat it as already imported.
			result.skipped++
			fmt.Printf("  skip  %-40s already stored\n", filename)
		case err != nil:
			result.failures = append(result.failures, importFailure{filename, err})
			fmt.Printf("  FAIL  %-40s %v\n", filename, err)
		default:
			result.added++
			fmt.Printf("  add   %-40s %6.1fs  %s\n", filename, seconds, name)
		}
	}

	if err := tx.Commit(); err != nil {
		cleanupCopies()
		return result, fmt.Errorf("commit import: %w", err)
	}
	return result, nil
}

// reportDryRun prints the outcome an import would have, without touching disk or the
// database. Durations are still measured, so this also surfaces unreadable files.
func reportDryRun(srcDir string, filenames []string, stored map[string]bool, overrides map[string]string) error {
	fmt.Printf("dry run: %d mp3 file(s) in %s\n\n", len(filenames), srcDir)

	var wouldAdd, wouldSkip, unreadable int
	for _, filename := range filenames {
		if stored[filename] {
			wouldSkip++
			fmt.Printf("  skip  %-40s already stored\n", filename)
			continue
		}

		seconds, err := audio.MP3Duration(filepath.Join(srcDir, filename))
		if err != nil {
			unreadable++
			fmt.Printf("  FAIL  %-40s %v\n", filename, err)
			continue
		}

		wouldAdd++
		fmt.Printf("  add   %-40s %6.1fs  %s\n", filename, seconds, nameFor(filename, overrides))
	}

	fmt.Printf("\nwould add %d, would skip %d (already stored), unreadable %d\n",
		wouldAdd, wouldSkip, unreadable)
	return nil
}

// storedFilenames returns the set of filenames already in the database, so an import can
// report skips up front instead of discovering each one via a failed insert.
func storedFilenames(ctx context.Context, queries *gen.Queries) (map[string]bool, error) {
	soundbites, err := queries.ListSoundbites(ctx)
	if err != nil {
		return nil, fmt.Errorf("list existing soundbites: %w", err)
	}

	stored := make(map[string]bool, len(soundbites))
	for _, sb := range soundbites {
		stored[sb.Filename] = true
	}
	return stored, nil
}
