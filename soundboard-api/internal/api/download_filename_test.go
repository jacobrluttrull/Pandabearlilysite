package api

import "testing"

func TestDownloadFilename(t *testing.T) {
	cases := []struct {
		name, display, stored, want string
	}{
		{"round-trips a well-named clip", "hi hi yes you AAAAA", "hi-hi-yes-you-AAAAA.mp3", "hi-hi-yes-you-AAAAA.mp3"},
		{"beats a run-together original", "ass eaten by these bitches", "asseatenbythesebitches.mp3", "ass-eaten-by-these-bitches.mp3"},
		{"preserves shouting", "IM DOIN TRICKS ON IT", "IMDOINTRICKSONIT.mp3", "IM-DOIN-TRICKS-ON-IT.mp3"},
		{"keeps variant markers", "im in your walls (+15db)", "x.mp3", "im-in-your-walls-(+15db).mp3"},
		{"no double hyphens", "a  --  b", "x.mp3", "a-b.mp3"},
		{"strips windows-reserved chars", `bad/name\with:chars`, "x.mp3", "badnamewithchars.mp3"},
		{"strips trailing dots", "trailing dots...", "x.mp3", "trailing-dots.mp3"},
		{"trims stray hyphens", "  -spaced-  ", "x.mp3", "spaced.mp3"},
		{"falls back when empty", "", "fallback.mp3", "fallback.mp3"},
		{"falls back when nothing survives", `<>:"/\|?*`, "fallback.mp3", "fallback.mp3"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := downloadFilename(c.display, c.stored); got != c.want {
				t.Errorf("downloadFilename(%q) = %q, want %q", c.display, got, c.want)
			}
		})
	}
}

func TestDownloadFilenameHasNoSpaces(t *testing.T) {
	for _, display := range []string{
		"hi hi yes you AAAAA", "im a gimme it now type a person", "you always see something in the rear",
	} {
		got := downloadFilename(display, "x.mp3")
		for _, r := range got {
			if r == ' ' {
				t.Errorf("downloadFilename(%q) = %q contains a space", display, got)
				break
			}
		}
	}
}

func TestDownloadDispositionIsAttachment(t *testing.T) {
	got := downloadDisposition("hi hi yes you AAAAA", "hi-hi-yes-you-AAAAA.mp3")
	want := `attachment; filename="hi-hi-yes-you-AAAAA.mp3"; filename*=UTF-8''hi-hi-yes-you-AAAAA.mp3`
	if got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}
