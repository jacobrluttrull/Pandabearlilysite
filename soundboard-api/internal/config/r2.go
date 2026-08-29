package config

import "soundboard-api/internal/clipstore"

// R2 gathers the Cloudflare R2 settings into the shape clipstore wants.
//
// Kept as its own method rather than having clipstore read the environment itself, so
// there stays exactly one place in the program that knows which variable names exist.
func (c Config) R2() clipstore.R2Config {
	return clipstore.R2Config{
		AccountID:       c.R2AccountID,
		AccessKeyID:     c.R2AccessKeyID,
		SecretAccessKey: c.R2SecretAccessKey,
		Bucket:          c.R2Bucket,
	}
}
