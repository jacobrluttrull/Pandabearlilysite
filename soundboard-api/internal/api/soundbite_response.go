package api

import (
	"fmt"

	"soundboard-api/internal/db/gen"
)

// soundbiteResponse is the JSON shape returned to the frontend. It hides
// storage details (filename) behind a playable audio_url.
type soundbiteResponse struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	DateMade      *string `json:"date_made"`
	DateStored    string  `json:"date_stored"`
	LengthSeconds float64 `json:"length_seconds"`
	AudioURL      string  `json:"audio_url"`
}

func toSoundbiteResponse(s gen.Soundbite) soundbiteResponse {
	resp := soundbiteResponse{
		ID:            s.ID,
		Name:          s.Name,
		DateStored:    s.DateStored,
		LengthSeconds: s.LengthSeconds,
		AudioURL:      fmt.Sprintf("/soundbites/%d/audio", s.ID),
	}
	if s.DateMade.Valid {
		resp.DateMade = &s.DateMade.String
	}
	return resp
}
