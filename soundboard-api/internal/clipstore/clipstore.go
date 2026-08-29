// Package clipstore holds clip audio and hands it to HTTP handlers.
//
// Two implementations exist because the same binary runs in two places. Locally the
// clips are files in a folder, which is what development wants: no credentials, no
// network, edit a file and it is served. In production they live in Cloudflare R2, which
// is what keeps 6 MB (and growing) of audio out of the git repository and out of the
// container image.
//
// The interface is deliberately about *serving* rather than *reading bytes*, because the
// two stores answer a request in fundamentally different ways: the local one streams the
// file, and the R2 one redirects the browser to fetch it directly. Modelling it as an
// io.Reader would have forced every byte through this process, which is the one thing
// object storage exists to avoid.
package clipstore

import (
	"context"
	"io"
	"net/http"
)

// Describer names a store in a form safe to log. Both the read and write sides need it,
// because both the API and the CLI state at startup which store they are talking to —
// the wrong one is otherwise invisible until a clip 404s.
type Describer interface {
	Describe() string
}

// Store serves clip audio over HTTP.
type Store interface {
	Describer

	// Serve responds with the audio for filename.
	//
	// downloadName, when non-empty, asks the browser to save the clip under that name
	// instead of playing it inline.
	Serve(w http.ResponseWriter, r *http.Request, filename, downloadName string)
}

// Manager is the write side, used by the CLI rather than the API. It is a separate
// interface because the running service never needs it: the API only ever reads, so
// giving it upload and delete would be handing it authority it has no use for.
type Manager interface {
	Describer

	// Put stores data under filename, overwriting any existing object.
	Put(ctx context.Context, filename string, data io.Reader, size int64) error

	// List returns every stored filename.
	List(ctx context.Context) ([]string, error)

	// Delete removes filename. Deleting something absent is not an error, so a repeated
	// cleanup does not fail the second time.
	Delete(ctx context.Context, filename string) error
}
