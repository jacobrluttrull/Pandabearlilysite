package audio

import (
	"errors"
	"io"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/tcolgate/mp3"
)

// Properties of the silent frames used to build test files: MPEG1 Layer3 at 44.1kHz
// carries 1152 samples per frame, which is 626 bytes on the wire.
const (
	silentFrameBytes   = 626
	samplesPerFrame    = 1152
	testFileSampleRate = 44100
)

// writeSilentMP3 writes an MP3 of exactly frameCount silent frames, optionally prefixed
// with leading bytes (used to stand in for an ID3 tag), and returns its path.
func writeSilentMP3(t *testing.T, frameCount int, prefix []byte) string {
	t.Helper()

	audio := make([]byte, frameCount*silentFrameBytes)
	if _, err := io.ReadFull(mp3.MakeSilence(), audio); err != nil {
		t.Fatalf("generate silence: %v", err)
	}

	path := filepath.Join(t.TempDir(), "clip.mp3")
	if err := os.WriteFile(path, append(prefix, audio...), 0o644); err != nil {
		t.Fatalf("write test mp3: %v", err)
	}
	return path
}

// id3v2Tag builds a minimal ID3v2 header with a zeroed payload, the kind of non-audio
// data real clips carry ahead of their first frame.
func id3v2Tag(payloadSize int) []byte {
	tag := []byte{'I', 'D', '3', 3, 0, 0}
	// Size is stored as four synchsafe bytes (7 significant bits each).
	for shift := 21; shift >= 0; shift -= 7 {
		tag = append(tag, byte((payloadSize>>shift)&0x7f))
	}
	return append(tag, make([]byte, payloadSize)...)
}

func TestMP3DurationMatchesSampleCount(t *testing.T) {
	for _, frameCount := range []int{1, 10, 383} {
		path := writeSilentMP3(t, frameCount, nil)

		got, err := MP3Duration(path)
		if err != nil {
			t.Fatalf("%d frames: MP3Duration: %v", frameCount, err)
		}

		// Derived from the format rather than from Frame.Duration, so this checks the
		// accumulation independently of the library's own arithmetic.
		want := float64(frameCount*samplesPerFrame) / testFileSampleRate
		if math.Abs(got-want) > 0.001 {
			t.Errorf("%d frames: got %.4fs, want %.4fs", frameCount, got, want)
		}
	}
}

func TestMP3DurationSkipsID3Tag(t *testing.T) {
	const frameCount = 20
	path := writeSilentMP3(t, frameCount, id3v2Tag(512))

	got, err := MP3Duration(path)
	if err != nil {
		t.Fatalf("MP3Duration: %v", err)
	}

	want := float64(frameCount*samplesPerFrame) / testFileSampleRate
	if math.Abs(got-want) > 0.001 {
		t.Errorf("got %.4fs, want %.4fs (tag bytes counted as audio?)", got, want)
	}
}

func TestMP3DurationToleratesTruncatedFinalFrame(t *testing.T) {
	const frameCount = 5
	path := writeSilentMP3(t, frameCount, nil)

	full, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read test mp3: %v", err)
	}
	// Lop off half of the last frame, as an editor trimming a clip might.
	truncated := filepath.Join(t.TempDir(), "truncated.mp3")
	if err := os.WriteFile(truncated, full[:len(full)-silentFrameBytes/2], 0o644); err != nil {
		t.Fatalf("write truncated mp3: %v", err)
	}

	got, err := MP3Duration(truncated)
	if err != nil {
		t.Fatalf("MP3Duration: %v", err)
	}

	// The four intact frames must still be measured.
	want := float64((frameCount-1)*samplesPerFrame) / testFileSampleRate
	if math.Abs(got-want) > 0.001 {
		t.Errorf("got %.4fs, want %.4fs", got, want)
	}
}

func TestMP3DurationRejectsNonMP3(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notaudio.mp3")
	if err := os.WriteFile(path, []byte("this is not an mp3 file"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if _, err := MP3Duration(path); !errors.Is(err, ErrNoFrames) {
		t.Errorf("got err %v, want ErrNoFrames", err)
	}
}
