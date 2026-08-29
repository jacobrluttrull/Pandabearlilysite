package clipstore

import (
	"fmt"
	"sort"
	"strings"
)

// RequireComplete returns an error if cfg is partly filled in.
//
// All four values or none is the only sensible state. A half-configured bucket almost
// always means a typo or a variable that failed to reach the process, and the fallback
// then serves clips from a directory that does not exist inside the container — every
// tile on the site 404s while the service reports itself healthy. Saying so at boot beats
// discovering it by clicking.
func RequireComplete(cfg R2Config) error {
	values := map[string]string{
		"R2_ACCOUNT_ID":        cfg.AccountID,
		"R2_ACCESS_KEY_ID":     cfg.AccessKeyID,
		"R2_SECRET_ACCESS_KEY": cfg.SecretAccessKey,
		"R2_BUCKET":            cfg.Bucket,
	}

	var set, unset []string
	for name, value := range values {
		if value == "" {
			unset = append(unset, name)
		} else {
			set = append(set, name)
		}
	}

	// All four, or none at all: both are deliberate.
	if len(set) == 0 || len(unset) == 0 {
		return nil
	}

	sort.Strings(set)
	sort.Strings(unset)
	return fmt.Errorf(
		"R2 is only partly configured: %s set, %s missing. Set all four or none — "+
			"a partial config silently falls back to local clip files, which do not exist in production",
		strings.Join(set, ", "), strings.Join(unset, ", "),
	)
}
