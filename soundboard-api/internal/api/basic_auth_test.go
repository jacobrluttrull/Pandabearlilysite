package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// okHandler stands in for the real router: if a request reaches it, the middleware let
// the request through.
func okHandler() (http.Handler, *bool) {
	reached := false
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	}), &reached
}

func TestWithBasicAuth(t *testing.T) {
	cases := []struct {
		name        string
		password    string
		path        string
		user, pass  string
		sendCreds   bool
		wantStatus  int
		wantReached bool
	}{
		{"no password configured lets everything through", "", "/", "", "", false, http.StatusOK, true},
		{"correct credentials", "hunter2", "/", "panda", "hunter2", true, http.StatusOK, true},
		{"wrong password", "hunter2", "/", "panda", "nope", true, http.StatusUnauthorized, false},
		{"wrong username", "hunter2", "/", "someone", "hunter2", true, http.StatusUnauthorized, false},
		{"no credentials at all", "hunter2", "/", "", "", false, http.StatusUnauthorized, false},
		// The Dockerfile HEALTHCHECK and Railway both probe this without credentials. If
		// auth covered it, every deploy would be reported unhealthy and rolled back.
		{"health check is exempt", "hunter2", healthPath, "", "", false, http.StatusOK, true},
		// The API is behind the same prompt as the pages: a soundbite list is not public
		// just because it is JSON.
		{"api routes are covered too", "hunter2", "/api/soundbites", "", "", false, http.StatusUnauthorized, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			next, reached := okHandler()
			handler := withBasicAuth("panda", c.password, next)

			req := httptest.NewRequest(http.MethodGet, c.path, nil)
			if c.sendCreds {
				req.SetBasicAuth(c.user, c.pass)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != c.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, c.wantStatus)
			}
			if *reached != c.wantReached {
				t.Errorf("handler reached = %v, want %v", *reached, c.wantReached)
			}
			if c.wantStatus == http.StatusUnauthorized && rec.Header().Get("WWW-Authenticate") == "" {
				t.Error("401 without WWW-Authenticate: the browser shows no password prompt")
			}
		})
	}
}

func TestWithNoIndex(t *testing.T) {
	next, _ := okHandler()
	rec := httptest.NewRecorder()
	withNoIndex(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if got := rec.Header().Get("X-Robots-Tag"); got != "noindex, nofollow, noarchive" {
		t.Errorf("X-Robots-Tag = %q, want the noindex set", got)
	}
}
