package mongodb

import (
	"kenja2/keywords"
	"testing"
)

func TestMongoKeywords(t *testing.T) {
	s := "$.,{}[]():;/`"
	clean, _ := keywords.CleanKeywords(__mongoReplacer, s)
	if clean != "" {
		panic("replacer fail")
	}

	s = "$.,{}[]():;/`$.,{}[]():;/`"
	clean, _ = keywords.CleanKeywords(__mongoReplacer, s)
	if clean != "" {
		panic("replacer fail")
	}

	s = "  a  b    c         de "
	clean, _ = keywords.CleanKeywords(__mongoReplacer, s)
	if clean != "a b c de" {
		panic("replacer fail")
	}

	s = "\ta\nb \vc\fd\r\ne"
	clean, _ = keywords.CleanKeywords(__mongoReplacer, s)
	if clean != "a b c d e" {
		panic("replacer fail")
	}

	s = "school music band club"
	clean, _ = keywords.CleanKeywords(__mongoReplacer, s)
	if clean != "school music band club" {
		panic("breaking keywords")
	}
}
