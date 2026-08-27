package api

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// staticHandler serves the prerendered frontend from dir.
//
// SvelteKit's static output writes each route as its own file, so a request for "/about"
// has to resolve to "about.html" — net/http would 404 it. Anything still unresolved falls
// back to the site's own 404 page.
//
// Files are served with http.ServeFile rather than http.FileServer because FileServer
// redirects any path ending in "index.html" back to its directory, which turns a request
// for "/" into a 301 instead of a page.
func staticHandler(dir string) http.Handler {
	root := filepath.Clean(dir)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// path.Clean, not filepath.Clean: a URL path is always slash-separated, and on
		// Windows filepath would mangle it into backslashes — and read a leading "//"
		// as a UNC share, leaving a slash behind that breaks the prefix check below.
		requested := strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(r.URL.Path, "/")), "/")

		// Hashed build assets never change under a given name, so they can be cached
		// hard. Everything else must revalidate or a deploy goes unnoticed.
		if strings.HasPrefix(requested, "_app/immutable/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		candidates := []string{requested}
		if requested == "" {
			candidates = []string{"index.html"}
		} else {
			candidates = append(candidates, requested+".html", requested+"/index.html")
		}

		for _, candidate := range candidates {
			if path, ok := existingFile(root, candidate); ok {
				http.ServeFile(w, r, path)
				return
			}
		}

		// Unknown path: the site's own 404 page, with the right status code.
		if path, ok := existingFile(root, "404.html"); ok {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, path)
			return
		}
		http.NotFound(w, r)
	})
}

// existingFile resolves name inside root, refusing anything that escapes it or is not a
// regular file.
func existingFile(root, name string) (string, bool) {
	path := filepath.Join(root, filepath.FromSlash(name))

	// filepath.Join cleans the result, so this catches "../" traversal attempts.
	if path != root && !strings.HasPrefix(path, root+string(os.PathSeparator)) {
		return "", false
	}

	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return "", false
	}
	return path, true
}
