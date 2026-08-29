package clipstore

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalRoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := NewLocal(dir)
	ctx := context.Background()

	want := []byte("not really an mp3, but bytes are bytes")
	if err := store.Put(ctx, "clip.mp3", bytes.NewReader(want), int64(len(want))); err != nil {
		t.Fatalf("Put: %v", err)
	}

	names, err := store.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 1 || names[0] != "clip.mp3" {
		t.Fatalf("List = %v, want [clip.mp3]", names)
	}

	rec := httptest.NewRecorder()
	store.Serve(rec, httptest.NewRequest(http.MethodGet, "/a", nil), "clip.mp3", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("Serve status = %d, want 200", rec.Code)
	}
	if !bytes.Equal(rec.Body.Bytes(), want) {
		t.Error("Serve returned different bytes than were stored")
	}

	if err := store.Delete(ctx, "clip.mp3"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	// Deleting something already gone must not error, so a repeated cleanup or a retry
	// after a partial failure does not fail the second time.
	if err := store.Delete(ctx, "clip.mp3"); err != nil {
		t.Errorf("Delete of a missing clip should be a no-op, got %v", err)
	}
}

func TestLocalServeMissing(t *testing.T) {
	rec := httptest.NewRecorder()
	NewLocal(t.TempDir()).Serve(rec, httptest.NewRequest(http.MethodGet, "/a", nil), "nope.mp3", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

// Filenames come from database rows. A row whose filename contained traversal must not
// be able to read or write outside the clips directory.
func TestLocalRefusesTraversal(t *testing.T) {
	dir := t.TempDir()
	outside := filepath.Join(filepath.Dir(dir), "escaped.mp3")
	store := NewLocal(dir)
	ctx := context.Background()

	if err := store.Put(ctx, "../escaped.mp3", bytes.NewReader([]byte("x")), 1); err == nil {
		t.Error("Put with traversal should be refused")
	}
	if _, err := os.Stat(outside); err == nil {
		t.Error("traversal wrote a file outside the clips directory")
	}

	rec := httptest.NewRecorder()
	store.Serve(rec, httptest.NewRequest(http.MethodGet, "/a", nil), "../escaped.mp3", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("Serve with traversal = %d, want 404", rec.Code)
	}
}

func TestLocalListMissingDir(t *testing.T) {
	names, err := NewLocal(filepath.Join(t.TempDir(), "absent")).List(context.Background())
	if err != nil {
		t.Fatalf("List of a missing dir should not error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("List = %v, want empty", names)
	}
}
