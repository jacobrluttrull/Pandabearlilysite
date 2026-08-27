package config

// Config holds runtime settings for the API and CLI, sourced from environment variables.
type Config struct {
	// DBPath is the filesystem path to the SQLite database file. In production this
	// points at a mounted volume, since it is the only state that must outlive a deploy.
	DBPath string
	// AudioDir is the folder holding clip audio. It ships inside the container image
	// rather than living on the volume, so it is effectively read-only in production.
	AudioDir string
	// Addr is the address the API listens on, e.g. ":8080".
	Addr string
	// AllowedOrigin is the value sent in Access-Control-Allow-Origin. Empty disables
	// CORS headers entirely, which is correct when the site is served from this same
	// process and no cross-origin request ever happens.
	AllowedOrigin string
	// StaticDir is the folder holding the built frontend. Empty serves the API alone.
	StaticDir string
}

// Load reads configuration from environment variables, applying local-dev defaults
// for anything unset.
func Load() Config {
	return Config{
		DBPath:        envOr("SOUNDBOARD_DB_PATH", "data/soundboard.db"),
		AudioDir:      envOr("SOUNDBOARD_AUDIO_DIR", "clips"),
		Addr:          listenAddr(),
		// Defaults to no CORS headers at all: the site is served by this same process,
		// so nothing is cross-origin. Set it explicitly only when the frontend is
		// hosted somewhere else.
		AllowedOrigin: envOr("SOUNDBOARD_ALLOWED_ORIGIN", ""),
		StaticDir:     envOr("SOUNDBOARD_STATIC_DIR", ""),
	}
}

// listenAddr resolves the listen address, preferring PORT.
//
// Railway, Fly, Heroku, and Cloud Run all inject PORT and health-check the port they
// assigned; a service that ignores it binds the wrong port and fails to deploy.
// SOUNDBOARD_ADDR stays available for local use and for setting a specific interface.
func listenAddr() string {
	if port := envOr("PORT", ""); port != "" {
		return ":" + port
	}
	return envOr("SOUNDBOARD_ADDR", ":8080")
}
