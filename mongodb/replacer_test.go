package mongodb

import (
	"kenja2/keywords"
	"testing"
)

func TestMongoKeywords(t *testing.T) {
	s := "$.,{}[]():;/`"
	clean := keywords.CleanKeywords(__MONGO_REPLACER, s)
	if clean != "" {
		panic("replacer fail")
	}

	s = `"'`
	clean = keywords.CleanKeywords(__MONGO_REPLACER, s)
	if clean != "" {
		panic("replacer fail")
	}

	s = "  a  b    c         de "
	clean = keywords.CleanKeywords(__MONGO_REPLACER, s)
	if clean != "a b c de" {
		panic("replacer fail")
	}

	s = "\ta\nb \vc\fd\r\ne"
	clean = keywords.CleanKeywords(__MONGO_REPLACER, s)
	if clean != "a b c d e" {
		panic("replacer fail")
	}

	s = "school music band club"
	clean = keywords.CleanKeywords(__MONGO_REPLACER, s)
	if clean != "school music band club" {
		panic("breaking keywords")
	}
}
