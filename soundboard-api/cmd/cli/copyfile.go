package main

import (
	"fmt"
	"io"
	"os"
)

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
