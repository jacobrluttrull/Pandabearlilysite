package api

import "net/http"

// withNoIndex asks crawlers not to index anything served here.
//
// robots.txt in the frontend's static files says the same thing, but it is only a request
// a crawler may honour or ignore, and it is not sent with the response at all. This header
// travels with every page and asset, so a crawler that skipped robots.txt still sees it.
//
// Unconditional on purpose: this site is private by design, so there is no configuration
// under which it should want to be indexed.
func withNoIndex(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive")
		next.ServeHTTP(w, r)
	})
}
