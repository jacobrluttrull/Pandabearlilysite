package config

// Config holds runtime settings for the API and CLI, sourced from environment variables.
type Config struct {
	// DBPath is the filesystem path to the SQLite database file.
	DBPath string
	// AudioDir is the folder where uploaded audio clips are stored.
	AudioDir string
	// Addr is the address the API listens on, e.g. ":8080".
	Addr string
	// AllowedOrigin is the value sent in Access-Control-Allow-Origin.
	AllowedOrigin string
}

// Load reads configuration from environment variables, applying local-dev defaults
// for anything unset.
func Load() Config {
	return Config{
		DBPath:        envOr("SOUNDBOARD_DB_PATH", "data/soundboard.db"),
		AudioDir:      envOr("SOUNDBOARD_AUDIO_DIR", "data/audio"),
		Addr:          envOr("SOUNDBOARD_ADDR", ":8080"),
		AllowedOrigin: envOr("SOUNDBOARD_ALLOWED_ORIGIN", "*"),
	}
}
