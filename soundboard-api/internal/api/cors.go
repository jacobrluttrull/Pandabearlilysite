package api

import "net/http"

// withCORS adds cross-origin headers when allowedOrigin is set.
//
// When the frontend is served by this same process the site and API share an origin, no
// cross-origin request ever happens, and the headers are omitted entirely rather than
// advertising access that nothing needs.
func withCORS(allowedOrigin string, next http.Handler) http.Handler {
	if allowedOrigin == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
