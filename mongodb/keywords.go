package mongodb

import (
	"slices"
	"strings"
	"unicode"
)

type Keywords string

var __PLACE_HOLDER = ""
var __fORBIDDEN = strings.NewReplacer(
	"$", __PLACE_HOLDER,
	".", __PLACE_HOLDER,
	",", __PLACE_HOLDER,
	"{", __PLACE_HOLDER,
	"}", __PLACE_HOLDER,
	"[", __PLACE_HOLDER,
	"]", __PLACE_HOLDER,
	"(", __PLACE_HOLDER,
	")", __PLACE_HOLDER,
	":", __PLACE_HOLDER,
	";", __PLACE_HOLDER,
	"/", __PLACE_HOLDER,
	"`", __PLACE_HOLDER,
	"'", __PLACE_HOLDER,
	`"`, __PLACE_HOLDER,
)

func spaceIsSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return ' '
	} else {
		return r
	}
}

func CleanKeywords(raw string) Keywords {
	s := __fORBIDDEN.Replace(raw)
	s = strings.Map(spaceIsSpace, s)
	sep := strings.Split(s, " ")
	sep = slices.DeleteFunc(sep, func(s string) bool {
		return len(strings.TrimSpace(s)) == 0
	})
	return Keywords(strings.Join(sep, " "))
}
