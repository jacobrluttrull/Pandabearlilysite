package api

import "net/http"

// handleSoundbiteAudio serves a clip for inline playback — no Content-Disposition, so the
// browser plays it rather than saving it. See handleSoundbiteDownload for the save path.
func (s *Server) handleSoundbiteAudio(w http.ResponseWriter, r *http.Request) {
	s.serveSoundbiteFile(w, r, nil)
}
