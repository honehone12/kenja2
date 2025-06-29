package keywords

import (
	"slices"
	"strings"
	"unicode"
)

type Keywords string

const PLACE_HOLDER = ""

func spaceIsSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return ' '
	} else {
		return r
	}
}

func CleanKeywords(replacer *strings.Replacer, raw string) Keywords {
	s := replacer.Replace(raw)
	s = strings.Map(spaceIsSpace, s)
	sep := strings.Split(s, " ")
	sep = slices.DeleteFunc(sep, func(s string) bool {
		return len(strings.TrimSpace(s)) == 0
	})
	return Keywords(strings.Join(sep, " "))
}
