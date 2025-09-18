package keywords

import (
	"slices"
	"strings"
	"unicode"
)

type Keywords string

func spaceIsSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return ' '
	} else {
		return r
	}
}

// returns cleaned keywords and the number of single character token
func CleanKeywords(replacer *strings.Replacer, raw string) (Keywords, int) {
	s := replacer.Replace(raw)
	s = strings.Map(spaceIsSpace, s)
	sep := strings.Split(s, " ")
	sep = slices.DeleteFunc(sep, func(s string) bool {
		return len(strings.TrimSpace(s)) == 0
	})

	singleTkns := 0
	for _, s := range sep {
		if len(s) == 1 {
			singleTkns++
		}
	}

	return Keywords(strings.Join(sep, " ")), singleTkns
}
