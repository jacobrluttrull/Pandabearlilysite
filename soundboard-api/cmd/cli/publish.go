package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
)

// publishClip stores srcPath's bytes in the clip store under filename.
//
// The file is read fully into memory first rather than streamed. Clips are a few hundred
// kilobytes at most, and it removes a real hazard: when the store is the local directory,
// streaming a source file into a Put that writes that same path would truncate the file
// while it was still being read. Buffering means the read finishes before the write starts.
func publishClip(ctx context.Context, st *store, srcPath, filename string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", srcPath, err)
	}

	if err := st.clips.Put(ctx, filename, bytes.NewReader(data), int64(len(data))); err != nil {
		return fmt.Errorf("publish %s: %w", filename, err)
	}
	return nil
}
