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

// NewHandler builds the full HTTP handler for the soundboard API.
func NewHandler(queries *gen.Queries, audioDir, allowedOrigin string) http.Handler {
	s := &Server{
		queries:  queries,
		audioDir: audioDir,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /soundbites", s.handleListSoundbites)
	mux.HandleFunc("GET /soundbites/{id}/audio", s.handleSoundbiteAudio)
	mux.HandleFunc("GET /soundbites/{id}/download", s.handleSoundbiteDownload)
	mux.HandleFunc("POST /soundbites/{id}/play", s.handlePlaySoundbite)

	return withCORS(allowedOrigin, mux)
}
