package api

import (
	"fmt"
	"strings"
	"unicode"
)

// downloadDisposition builds the Content-Disposition header that saves a clip under a
// hyphenated form of its display name — "ass-eaten-by-these-bitches.mp3" rather than the
// storage filename "asseatenbythesebitches.mp3".
//
// The plain filename parameter is restricted to ASCII for older clients, and filename*
// carries the full UTF-8 name per RFC 5987 for everything current. Browsers prefer
// filename* when they understand it.
func downloadDisposition(displayName, storedFilename string) string {
	pretty := downloadFilename(displayName, storedFilename)

	ascii := toASCIIFilename(pretty)
	if ascii == "" {
		ascii = storedFilename
	}

	return fmt.Sprintf("attachment; filename=%q; filename*=UTF-8''%s",
		ascii, escapeExtValue(pretty))
}

// downloadFilename turns a display name into a filename that is safe on Windows, macOS,
// and Linux, falling back to the stored filename if nothing usable survives.
//
// Spaces become hyphens. A name with spaces is awkward to type, quote in a shell, or
// read in a downloads folder, and for a well-named clip the hyphenated form lands back on
// exactly the original filename ("hi hi yes you AAAAA" -> hi-hi-yes-you-AAAAA.mp3). For a
// clip whose filename was an unreadable run-together blob, it beats the original outright
// (asseatenbythesebitches.mp3 -> ass-eaten-by-these-bitches.mp3).
//
// Case is preserved, so IM-DOIN-TRICKS-ON-IT stays distinct from im-doing-tricks-on-it.
func downloadFilename(displayName, storedFilename string) string {
	cleaned := strings.Map(func(r rune) rune {
		// Reserved on Windows, plus path separators and control characters.
		if strings.ContainsRune(`<>:"/\|?*`, r) || unicode.IsControl(r) {
			return -1
		}
		return r
	}, displayName)

	// Collapse each run of whitespace into a single hyphen.
	cleaned = strings.Join(strings.Fields(cleaned), "-")

	// Collapse hyphen runs left by names that already contained one.
	for strings.Contains(cleaned, "--") {
		cleaned = strings.ReplaceAll(cleaned, "--", "-")
	}

	// Trailing dots and spaces are stripped by Windows, which would corrupt the name;
	// leading or trailing hyphens just look like a mistake.
	cleaned = strings.Trim(cleaned, "-. ")

	if cleaned == "" {
		return storedFilename
	}
	return cleaned + ".mp3"
}

// toASCIIFilename drops anything outside printable ASCII, since the bare filename
// parameter of Content-Disposition cannot carry more than that.
func toASCIIFilename(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 0x20 && r < 0x7f && r != '"' && r != '\\' {
			b.WriteRune(r)
		}
	}
	return strings.TrimRight(strings.Join(strings.Fields(b.String()), " "), ". ")
}

// escapeExtValue percent-encodes a name for the filename* parameter.
func escapeExtValue(s string) string {
	// url.PathEscape leaves through characters that some parsers choke on inside a
	// header parameter, so encode conservatively by hand against an explicit safe set.
	const safe = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$&+-.^_`|~"

	var b strings.Builder
	for _, c := range []byte(s) {
		if strings.IndexByte(safe, c) >= 0 {
			b.WriteByte(c)
			continue
		}
		fmt.Fprintf(&b, "%%%02X", c)
	}
	return b.String()
}
