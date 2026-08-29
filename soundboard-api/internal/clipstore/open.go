package clipstore

// Open picks the store to use: R2 when it is fully configured, the local directory
// otherwise.
//
// The fallback is silent by design, because a fresh checkout with no credentials must
// still run. Callers log Describe() at boot so the choice is visible rather than assumed
// — the same treatment the database gets, and for the same reason: the wrong one here
// looks healthy right up until a clip 404s.
func Open(cfg R2Config, localDir string) Store {
	if cfg.Configured() {
		return NewR2(cfg)
	}
	return NewLocal(localDir)
}
