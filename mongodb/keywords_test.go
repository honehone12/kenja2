package mongodb

import (
	"testing"
)

func TestKeywords(t *testing.T) {
	s := "$.,{}[]():;/`"
	clean := CleanKeywords(s)
	if clean != "" {
		panic("replacer fail")
	}

	s = `"'`
	clean = CleanKeywords(s)
	if clean != "" {
		panic("replacer fail")
	}

	s = "  a  b    c         de "
	clean = CleanKeywords(s)
	if clean != "a b c de" {
		panic("replacer fail")
	}

	s = "\ta\nb \vc\fd\r\ne"
	clean = CleanKeywords(s)
	if clean != "a b c d e" {
		panic("replacer fail")
	}

	s = "school music band club"
	clean = CleanKeywords(s)
	if clean != "school music band club" {
		panic("breaking keywords")
	}
}
