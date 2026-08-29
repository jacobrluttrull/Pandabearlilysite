package api

import (
	"net/http"

	"soundboard-api/internal/clipstore"
	"soundboard-api/internal/db/gen"
)

// Server holds the dependencies needed by route handlers.
type Server struct {
	queries *gen.Queries
	clips   clipstore.Store
}

// Options configures the HTTP handler.
type Options struct {
	// Clips is where clip audio comes from — a local directory or an R2 bucket.
	Clips clipstore.Store
	// AllowedOrigin sets Access-Control-Allow-Origin. Empty omits CORS headers, which
	// is right when the frontend is served by this same process.
	AllowedOrigin string
	// StaticDir is the built frontend to serve. Empty serves the API alone.
	StaticDir string
	// AuthUser and AuthPassword gate the whole site behind a password prompt. An empty
	// password disables the prompt.
	AuthUser     string
	AuthPassword string
}

// NewHandler builds the full HTTP handler for the soundboard.
//
// API routes live under /api so they cannot collide with a page route as the site grows;
// everything else falls through to the static frontend when one is configured.
func NewHandler(queries *gen.Queries, opts Options) http.Handler {
	s := &Server{
		queries: queries,
		clips:   opts.Clips,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/soundbites", s.handleListSoundbites)
	mux.HandleFunc("GET /api/soundbites/{id}/audio", s.handleSoundbiteAudio)
	mux.HandleFunc("GET /api/soundbites/{id}/download", s.handleSoundbiteDownload)
	mux.HandleFunc("POST /api/soundbites/{id}/play", s.handlePlaySoundbite)
	mux.HandleFunc("GET /api/health", handleHealth)

	if opts.StaticDir != "" {
		mux.Handle("/", staticHandler(opts.StaticDir))
	}

	// Outermost first: the no-index header travels with every response including a 401,
	// and the password check runs before anything routes, so an unauthenticated request
	// never reaches a handler or the database.
	return withNoIndex(withBasicAuth(opts.AuthUser, opts.AuthPassword, withCORS(opts.AllowedOrigin, mux)))
}

// handleHealth reports that the process is up. Platform health checks hit this rather
// than /api/soundbites so a check never depends on the database having rows.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
