package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
)

// playResponse is returned after recording a play, so the frontend can reconcile its
// optimistic count with the real total.
type playResponse struct {
	PlayCount int64 `json:"play_count"`
}

// handlePlaySoundbite records one play of a clip and returns the new total.
//
// Clips are only a second or two long, so a play is counted when playback starts rather
// than when it finishes — pressing the tile is the play.
func (s *Server) handlePlaySoundbite(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid soundbite id")
		return
	}

	playCount, err := s.queries.IncrementPlayCount(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "soundbite not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to record play")
		return
	}

	writeJSON(w, http.StatusOK, playResponse{PlayCount: playCount})
}
