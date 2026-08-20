package api

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (s *Server) handleSoundbiteAudio(w http.ResponseWriter, r *http.Request) {
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

	http.ServeContent(w, r, soundbite.Filename, info.ModTime(), file)
}
