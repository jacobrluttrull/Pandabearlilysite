package clipstore

import (
	"fmt"
	"os"
)

// RequireLocalDir returns an error if dir is not a directory that exists.
//
// This is the guard for the case RequireComplete cannot see: no R2 variables at all.
// A completely empty R2 config is legitimate — it is what a fresh checkout looks like —
// so it falls back to the local store rather than failing. In the container that fallback
// is a trap, because the image carries no clips directory since the audio moved to R2:
// every store call would miss, every tile on the site would 404, and the health check
// would keep reporting the service healthy because the API itself is fine.
//
// Checking the directory exists separates the two cases without the API having to know
// whether it is running in production. On a development machine the folder is there and
// this passes. In the image it is not, so a deploy that forgot the R2 variables dies at
// boot with the reason instead of serving a soundboard where nothing plays.
func RequireLocalDir(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf(
			"no clip store: R2 is unconfigured and the local clip directory %q does not exist. "+
				"Set R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY and R2_BUCKET, "+
				"or point SOUNDBOARD_AUDIO_DIR at a folder of clips", dir)
	}
	if err != nil {
		return fmt.Errorf("check local clip directory %q: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("local clip directory %q is a file, not a directory", dir)
	}
	return nil
}
