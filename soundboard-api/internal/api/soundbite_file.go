package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"soundboard-api/internal/db/gen"
)

// serveSoundbiteFile looks up a clip by the id in the request path and serves its audio.
//
// dispositionFor, when non-nil, is called with the looked-up clip to build the
// Content-Disposition value. The audio route passes nil so browsers play the clip inline;
// the download route returns an "attachment" disposition so the same bytes are saved
// instead. Passing a callback rather than a string keeps the lookup here, so the download
// route does not query the clip a second time just to name it.
//
// How the bytes actually reach the caller is the store's business: the local store
// streams the file with range support, and the R2 store redirects to a presigned URL so
// the audio never passes through this process at all.
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

	var disposition string
	if dispositionFor != nil {
		disposition = dispositionFor(soundbite)
	}

	s.clips.Serve(w, r, soundbite.Filename, disposition)
}
