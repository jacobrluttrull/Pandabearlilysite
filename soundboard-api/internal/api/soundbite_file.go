package api

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"soundboard-api/internal/db/gen"
)

// serveSoundbiteFile looks up a clip by the id in the request path and serves its audio.
//
// dispositionFor, when non-nil, is called with the looked-up clip to build the
// Content-Disposition header. The audio route passes nil so browsers play the clip
// inline; the download route returns an "attachment" disposition so the same bytes are
// saved to disk instead. Passing a callback rather than a string keeps the lookup here,
// so the download route does not have to query the clip a second time just to name it.
//
// http.ServeContent handles range requests, so seeking works and a partial download can
// resume.
func (s *Server) serveSoundbiteFile(
	w http.ResponseWriter,
	r *http.Request,
	dispositionFor func(gen.Soundbite) string,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid soundbite id")
		return
	}

	soundbite, err := s.queries.GetSoundbite(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "soundbite not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up soundbite")
		return
	}

	path := filepath.Join(s.audioDir, soundbite.Filename)

	file, err := os.Open(path)
	if err != nil {
		writeError(w, http.StatusNotFound, "audio file not found")
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read audio file")
		return
	}

	if dispositionFor != nil {
		w.Header().Set("Content-Disposition", dispositionFor(soundbite))
	}

	http.ServeContent(w, r, soundbite.Filename, info.ModTime(), file)
}
