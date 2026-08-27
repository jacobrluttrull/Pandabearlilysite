package api

import (
	"net/http"

	"soundboard-api/internal/db/gen"
)

// handleSoundbiteDownload serves a clip as a file to save rather than audio to play.
//
// This route exists because the HTML `download` attribute is ignored cross-origin, and
// the site is served from a different origin than this API. Without an attachment header
// the browser would simply navigate to the clip and play it.
//
// The saved file is named after the clip's display name, so a visitor gets
// "hi hi yes you AAAAA.mp3" rather than the storage filename.
func (s *Server) handleSoundbiteDownload(w http.ResponseWriter, r *http.Request) {
	s.serveSoundbiteFile(w, r, func(sb gen.Soundbite) string {
		return downloadDisposition(sb.Name, sb.Filename)
	})
}
