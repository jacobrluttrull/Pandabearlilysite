package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// buildDir lays out a directory shaped like SvelteKit's static output.
func buildDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	files := map[string]string{
		"index.html":                  "<title>home</title>",
		"about.html":                  "<title>about</title>",
		"404.html":                    "<title>not found</title>",
		"robots.txt":                  "User-agent: *",
		"images/pfp.png":              "\x89PNG",
		"_app/immutable/entry/app.js": "console.log(1)",
		"nested/index.html":           "<title>nested</title>",
		"secret.txt":                  "should not escape",
	}
	for name, body := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestStaticHandlerServesPages(t *testing.T) {
	h := staticHandler(buildDir(t))

	cases := []struct{ path, wantBody string }{
		{"/", "home"},            // root must serve index.html, not redirect
		{"/about", "about"},      // extensionless route -> about.html
		{"/about.html", "about"}, // explicit file still works
		{"/nested", "nested"},    // directory route -> nested/index.html
		{"/robots.txt", "User-agent"},
	}
	for _, c := range cases {
		rec := get(t, h, c.path)
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", c.path, rec.Code)
			continue
		}
		if body := rec.Body.String(); !contains(body, c.wantBody) {
			t.Errorf("GET %s body = %q, want it to contain %q", c.path, body, c.wantBody)
		}
	}
}

func TestStaticHandlerUnknownPathServes404Page(t *testing.T) {
	h := staticHandler(buildDir(t))

	rec := get(t, h, "/no-such-page")
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
	if !contains(rec.Body.String(), "not found") {
		t.Errorf("body = %q, want the site's 404 page", rec.Body.String())
	}
}

// The 404 page used to go out through http.ServeFile after the status had already been
// written, so ServeFile's headers were set too late to be sent and Go logged a
// superfluous WriteHeader call on every miss. The status was right either way, which is
// why this asserts the headers instead — they are what the bug actually cost.
func TestStaticHandler404PageSendsItsHeaders(t *testing.T) {
	h := staticHandler(buildDir(t))

	rec := get(t, h, "/no-such-page")
	if got := rec.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html; charset=utf-8", got)
	}
	if got, want := rec.Header().Get("Content-Length"), strconv.Itoa(rec.Body.Len()); got != want {
		t.Errorf("Content-Length = %q, want %q", got, want)
	}
}

func TestStaticHandlerCachesImmutableAssets(t *testing.T) {
	h := staticHandler(buildDir(t))

	rec := get(t, h, "/_app/immutable/entry/app.js")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Cache-Control"); got != "public, max-age=31536000, immutable" {
		t.Errorf("Cache-Control = %q, want the immutable directive", got)
	}
}

func TestStaticHandlerDoesNotCacheHTML(t *testing.T) {
	h := staticHandler(buildDir(t))

	if got := get(t, h, "/about").Header().Get("Cache-Control"); got != "" {
		t.Errorf("Cache-Control = %q, want none so a deploy is picked up", got)
	}
}

func TestStaticHandlerRefusesTraversal(t *testing.T) {
	dir := buildDir(t)
	// A file outside the served directory that must never be reachable.
	outside := filepath.Join(filepath.Dir(dir), "outside.txt")
	if err := os.WriteFile(outside, []byte("private"), 0o644); err != nil {
		t.Fatal(err)
	}

	h := staticHandler(dir)
	for _, path := range []string{"/../outside.txt", "/images/../../outside.txt", "/%2e%2e/outside.txt"} {
		rec := get(t, h, path)
		if contains(rec.Body.String(), "private") {
			t.Errorf("GET %s leaked a file outside the build dir", path)
		}
	}
}

func contains(haystack, needle string) bool {
	return len(needle) > 0 && len(haystack) >= len(needle) &&
		(func() bool {
			for i := 0; i+len(needle) <= len(haystack); i++ {
				if haystack[i:i+len(needle)] == needle {
					return true
				}
			}
			return false
		})()
}
