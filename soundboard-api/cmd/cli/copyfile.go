package main

import (
	"fmt"
	"io"
	"os"
)

// copyFileIfMissing copies srcPath to destPath unless destPath already exists, and
// reports whether it actually wrote anything. Bulk imports use this so a re-run after a
// partial failure skips clips already in the audio dir instead of erroring out.
func copyFileIfMissing(srcPath, destPath string) (copied bool, err error) {
	if _, err := os.Stat(destPath); err == nil {
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}

	if err := copyFile(srcPath, destPath); err != nil {
		return false, err
	}
	return true, nil
}

func copyFile(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("%s already exists", destPath)
	}

	dest, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		os.Remove(destPath)
		return err
	}

	return nil
}
