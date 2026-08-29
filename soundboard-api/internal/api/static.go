package api

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
		switch {
		case strings.HasPrefix(requested, "_app/immutable/"):
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

		// Images and fonts are the opposite case: heavy, but named by hand, so they
		// cannot be immutable. With no header at all a browser revalidates every one on
		// every visit — the avatar alone is 147 KB re-fetched on each page load. A week
		// of caching removes that; the cost is that replacing a file under the same name
		// takes up to a week to reach someone who has already seen it, so change the
		// filename when the picture itself changes.
		case isLongLivedAsset(requested):
			w.Header().Set("Cache-Control", "public, max-age=604800")
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
		//
		// Written out directly rather than with http.ServeFile, which always sends 200
		// and has no way to be told otherwise. Pairing it with an earlier WriteHeader(404)
		// does produce a 404 — the first status written wins — but ServeFile then tries to
		// write its own, which Go reports as "superfluous response.WriteHeader call" on
		// every miss, and the Content-Type and Content-Length it sets arrive too late to
		// be sent at all.
		if path, ok := existingFile(root, "404.html"); ok {
			if page, err := os.ReadFile(path); err == nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("Content-Length", strconv.Itoa(len(page)))
				w.WriteHeader(http.StatusNotFound)
				w.Write(page)
				return
			}
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

// longLivedAssetExts are file types that are large, static, and referenced by a name the
// site author chooses. HTML is deliberately absent: it is the file that names every other
// one, so caching it is what makes a deploy go unnoticed.
var longLivedAssetExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".avif": true,
	".gif": true, ".svg": true, ".ico": true,
	".woff": true, ".woff2": true, ".ttf": true, ".otf": true,
	".mp3": true, ".mp4": true, ".webm": true,
}

// isLongLivedAsset reports whether name is one of those file types.
func isLongLivedAsset(name string) bool {
	return longLivedAssetExts[strings.ToLower(path.Ext(name))]
}
