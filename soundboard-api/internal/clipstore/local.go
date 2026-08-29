package clipstore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Local serves clips from a directory on disk. This is the development store, and the
// one that runs when no R2 credentials are configured.
type Local struct {
	dir string
}

// NewLocal returns a store reading from dir.
func NewLocal(dir string) *Local {
	return &Local{dir: filepath.Clean(dir)}
}

func (l *Local) Describe() string { return "local directory " + l.dir }

// Serve streams the file. http.ServeContent handles range requests, so seeking works and
// a partial download can resume — the behaviour the R2 store gets from Cloudflare instead.
func (l *Local) Serve(w http.ResponseWriter, r *http.Request, filename, downloadName string) {
	path, ok := l.resolve(filename)
	if !ok {
		http.Error(w, "audio file not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "audio file not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, "failed to read audio file", http.StatusInternalServerError)
		return
	}

	if downloadName != "" {
		w.Header().Set("Content-Disposition", downloadName)
	}
	http.ServeContent(w, r, filename, info.ModTime(), file)
}

func (l *Local) Put(ctx context.Context, filename string, data io.Reader, size int64) error {
	path, ok := l.resolve(filename)
	if !ok {
		return fmt.Errorf("refusing to write outside %s: %q", l.dir, filename)
	}
	if err := os.MkdirAll(l.dir, 0o755); err != nil {
		return fmt.Errorf("create clips dir: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", filename, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, data); err != nil {
		return fmt.Errorf("write %s: %w", filename, err)
	}
	return nil
}

func (l *Local) List(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(l.dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read clips dir: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ".mp3") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func (l *Local) Delete(ctx context.Context, filename string) error {
	path, ok := l.resolve(filename)
	if !ok {
		return fmt.Errorf("refusing to delete outside %s: %q", l.dir, filename)
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete %s: %w", filename, err)
	}
	return nil
}

// resolve joins filename onto the store's directory, refusing anything that escapes it.
// Filenames come from the database, so this guards against a row whose filename contains
// traversal rather than against untrusted request input.
func (l *Local) resolve(filename string) (string, bool) {
	path := filepath.Join(l.dir, filepath.FromSlash(filename))
	if path != l.dir && !strings.HasPrefix(path, l.dir+string(os.PathSeparator)) {
		return "", false
	}
	return path, true
}
