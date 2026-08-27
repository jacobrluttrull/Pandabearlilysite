package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// defaultNamesPath is where curated display names live.
//
// It sits outside data/ deliberately: data/ is gitignored local runtime state, but these
// names are hand-written content. Losing them means retyping every clip label, so they
// belong in version control.
const defaultNamesPath = "names.json"

// loadNameOverrides reads the filename -> display name map used to correct clips whose
// filenames do not derive into anything readable.
//
// A missing file is not an error; it simply means nothing has been renamed yet.
func loadNameOverrides(path string) (map[string]string, error) {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read names file: %w", err)
	}

	overrides := map[string]string{}
	if err := json.Unmarshal(raw, &overrides); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return overrides, nil
}

// saveNameOverrides writes the map back out. Go marshals map keys in sorted order, so
// the file stays stable and diff-friendly across runs.
func saveNameOverrides(path string, overrides map[string]string) error {
	raw, err := json.MarshalIndent(overrides, "", "  ")
	if err != nil {
		return fmt.Errorf("encode names: %w", err)
	}

	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		return fmt.Errorf("write names file: %w", err)
	}
	return nil
}

// nameFor resolves the label for a clip: a curated name if one is listed, otherwise the
// name derived from the filename. Blank entries fall back to derivation, so emptying a
// value in the file is a way to undo an override rather than a way to blank a label.
func nameFor(filename string, overrides map[string]string) string {
	if custom, ok := overrides[filename]; ok && strings.TrimSpace(custom) != "" {
		return strings.TrimSpace(custom)
	}
	return displayName(filename)
}
