package api

import "net/http"

func (s *Server) handleListSoundbites(w http.ResponseWriter, r *http.Request) {
	soundbites, err := s.queries.ListSoundbites(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list soundbites")
		return
	}

	resp := make([]soundbiteResponse, len(soundbites))
	for i, sb := range soundbites {
		resp[i] = toSoundbiteResponse(sb)
	}

	writeJSON(w, http.StatusOK, resp)
}
