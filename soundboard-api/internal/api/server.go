package api

import (
	"net/http"

	"soundboard-api/internal/db/gen"
)

// Server holds the dependencies needed by route handlers.
type Server struct {
	queries  *gen.Queries
	audioDir string
}

// Options configures the HTTP handler.
type Options struct {
	// AudioDir is where clip audio is read from.
	AudioDir string
	// AllowedOrigin sets Access-Control-Allow-Origin. Empty omits CORS headers, which
	// is right when the frontend is served by this same process.
	AllowedOrigin string
	// StaticDir is the built frontend to serve. Empty serves the API alone.
	StaticDir string
}

// NewHandler builds the full HTTP handler for the soundboard.
//
// API routes live under /api so they cannot collide with a page route as the site grows;
// everything else falls through to the static frontend when one is configured.
func NewHandler(queries *gen.Queries, opts Options) http.Handler {
	s := &Server{
		queries:  queries,
		audioDir: opts.AudioDir,
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

	return withCORS(opts.AllowedOrigin, mux)
}

// handleHealth reports that the process is up. Platform health checks hit this rather
// than /api/soundbites so a check never depends on the database having rows.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
