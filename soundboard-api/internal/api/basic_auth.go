package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// healthPath is exempt from authentication. Railway and the Dockerfile HEALTHCHECK both
// probe it without credentials, so requiring a password here would report every deploy as
// unhealthy and roll it back. It reveals nothing but "the process is up".
const healthPath = "/api/health"

// withBasicAuth puts the whole site behind a password when one is configured.
//
// This is a fan site for someone else's character, published without her endorsement, so
// it is deliberately not a public website — the password is what keeps a link shareable
// without the site being open to the world. An empty password disables the check, which
// is what local development wants; the API logs which mode it is in at boot so an
// unprotected deploy is visible rather than assumed.
func withBasicAuth(username, password string, next http.Handler) http.Handler {
	if password == "" {
		return next
	}

	// Hashing both sides first means the comparison is over two fixed-length digests, so
	// it leaks neither the length of the real credentials nor, via early exit, how many
	// leading characters a guess got right.
	wantUser := sha256.Sum256([]byte(username))
	wantPass := sha256.Sum256([]byte(password))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == healthPath {
			next.ServeHTTP(w, r)
			return
		}

		gotUser, gotPass, ok := r.BasicAuth()
		if ok {
			haveUser := sha256.Sum256([]byte(gotUser))
			havePass := sha256.Sum256([]byte(gotPass))

			// Both compared every time — no && short-circuit — so a wrong username costs
			// the same as a wrong password and neither can be probed separately.
			userMatch := subtle.ConstantTimeCompare(haveUser[:], wantUser[:])
			passMatch := subtle.ConstantTimeCompare(havePass[:], wantPass[:])
			if userMatch&passMatch == 1 {
				next.ServeHTTP(w, r)
				return
			}
		}

		// The realm string is shown by the browser in its password prompt.
		w.Header().Set("WWW-Authenticate", `Basic realm="PandaLily fan site", charset="UTF-8"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}
