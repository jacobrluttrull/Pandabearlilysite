// Package clipname derives and stores the labels shown for each soundbite.
package clipname

import (
	"path/filepath"
	"strings"
	"unicode"
)

// DisplayName derives a clip's label from its filename, staying deliberately faithful to
// the file so you can read a name on the grid and know which file it came from:
//
//	ads-in-our-ass.mp3   -> "ads in our ass"
//	distorted_i_busted   -> "distorted i busted"
//	SheDidWhat.mp3       -> "She Did What"
//	hi-hi-yes-you-AAAAA  -> "hi hi yes you AAAAA"
//	IMDOINTRICKSONIT.mp3 -> "IMDOINTRICKSONIT"
//
// Casing is preserved rather than normalised: a name shouted in the filename is shouted
// on the grid, and IMDOINTRICKSONIT stays distinguishable from imdoingtricksonit.
//
// Filenames written as one run-together word cannot be split by any reliable rule, so
// they come through unchanged and are meant to be corrected in the names file. See
// LoadOverrides.
func DisplayName(filename string) string {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	base = splitSeparators(base)
	base = splitCamelCase(base)

	name := strings.Join(strings.Fields(base), " ")
	if name == "" {
		// A filename made entirely of separators still needs some label.
		return filename
	}
	return name
}

// splitSeparators turns the punctuation used to join words in filenames into spaces.
func splitSeparators(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '-' || r == '_' || r == '.' {
			return ' '
		}
		return r
	}, s)
}

// splitCamelCase inserts a space at each lower-to-upper boundary, which recovers the
// word breaks in names like "SheDidWhat" and "lilyYEAHHHH". Names in a single case, such
// as "IMDOINTRICKSONIT", have no boundary to find and are left alone.
func splitCamelCase(s string) string {
	var out strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			prev := rune(s[i-1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) {
				out.WriteRune(' ')
			}
		}
		out.WriteRune(r)
	}
	return out.String()
}
