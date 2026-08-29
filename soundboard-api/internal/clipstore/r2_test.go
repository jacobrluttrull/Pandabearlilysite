package clipstore

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func testR2() *R2 {
	return NewR2(R2Config{
		AccountID:       "acct",
		AccessKeyID:     "key",
		SecretAccessKey: "secret",
		Bucket:          "clips",
	})
}

func TestR2ConfigConfigured(t *testing.T) {
	full := R2Config{AccountID: "a", AccessKeyID: "k", SecretAccessKey: "s", Bucket: "b"}
	if !full.Configured() {
		t.Fatal("a complete config should be Configured")
	}

	// Every field is required: a half-set config must fall back to the local store
	// rather than build a client that fails on first use.
	for _, missing := range []string{"AccountID", "AccessKeyID", "SecretAccessKey", "Bucket"} {
		c := full
		switch missing {
		case "AccountID":
			c.AccountID = ""
		case "AccessKeyID":
			c.AccessKeyID = ""
		case "SecretAccessKey":
			c.SecretAccessKey = ""
		case "Bucket":
			c.Bucket = ""
		}
		if c.Configured() {
			t.Errorf("config missing %s should not be Configured", missing)
		}
	}
}

// Serving must redirect rather than stream: sending the bytes through this process is
// exactly what moving audio to object storage is meant to avoid.
func TestR2ServeRedirects(t *testing.T) {
	rec := httptest.NewRecorder()
	testR2().Serve(rec, httptest.NewRequest(http.MethodGet, "/api/soundbites/1/audio", nil), "clip.mp3", "")

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
	}

	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "acct.r2.cloudflarestorage.com") || !strings.Contains(loc, "clip.mp3") {
		t.Errorf("Location does not point at the object: %s", loc)
	}
	if !strings.Contains(loc, "X-Amz-Signature") {
		t.Error("Location is not a presigned URL — it would be rejected by a private bucket")
	}
	// A presigned URL expires and is minted per request, so caching it as a permanent
	// location would hand out dead links.
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", rec.Header().Get("Cache-Control"))
	}
}

// The download route names the saved file after the clip's display name. A redirect
// discards any header set here, so the disposition has to be signed into the URL instead
// — this is the test that would catch it silently reverting to the storage filename.
func TestR2ServeSignsDownloadName(t *testing.T) {
	rec := httptest.NewRecorder()
	disposition := `attachment; filename="ass-eaten-by-these-bitches.mp3"`
	testR2().Serve(rec, httptest.NewRequest(http.MethodGet, "/d", nil), "asseatenbythesebitches.mp3", disposition)

	loc, err := url.Parse(rec.Header().Get("Location"))
	if err != nil {
		t.Fatalf("Location is not a URL: %v", err)
	}
	got := loc.Query().Get("response-content-disposition")
	if got != disposition {
		t.Errorf("response-content-disposition = %q, want %q", got, disposition)
	}
	if !strings.Contains(loc.Query().Get("X-Amz-SignedHeaders")+loc.RawQuery, "X-Amz-Signature") {
		t.Error("disposition must be part of the signed URL, not an unsigned addition")
	}
}

func TestOpenPicksStore(t *testing.T) {
	full := R2Config{AccountID: "a", AccessKeyID: "k", SecretAccessKey: "s", Bucket: "b"}
	if _, ok := Open(full, "clips").(*R2); !ok {
		t.Error("a configured R2Config should select the R2 store")
	}
	if _, ok := Open(R2Config{}, "clips").(*Local); !ok {
		t.Error("an empty R2Config should fall back to the local store")
	}
}
