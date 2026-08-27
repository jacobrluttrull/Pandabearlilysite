// Package audio inspects local audio files for the metadata the soundboard stores.
package audio

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/tcolgate/mp3"
)

// ErrNoFrames is returned when a file contains no decodable MP3 frames, which usually
// means it is not really an MP3 despite its extension.
var ErrNoFrames = errors.New("no mp3 frames found")

// MP3Duration returns the playing length of an MP3 file in seconds.
//
// It sums the duration of every frame rather than estimating from file size and
// bitrate, so variable-bitrate clips measure correctly. Only frame headers are read;
// the audio itself is never decoded, which keeps a few hundred clips fast to scan.
// Leading non-audio bytes such as ID3 tags are skipped by the decoder.
func MP3Duration(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var (
		decoder = mp3.NewDecoder(file)
		frame   mp3.Frame
		skipped int
		seconds float64
		frames  int
	)

	for {
		err := decoder.Decode(&frame, &skipped)
		if err != nil {
			// A truncated final frame is common in clips trimmed by editors; as long as
			// we read real frames before it, the measured length is still good.
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			return 0, fmt.Errorf("decode frame %d: %w", frames+1, err)
		}
		seconds += frame.Duration().Seconds()
		frames++
	}

	if frames == 0 {
		return 0, ErrNoFrames
	}
	return seconds, nil
}
