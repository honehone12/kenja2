package opensearch

import (
	"kenja2/keywords"
	"testing"
)

func TestOpensearchKeywords(t *testing.T) {
	s := "_+-*%&|=<>!?^~{}[]():,./'`"
	clean := keywords.CleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "" {
		panic("replacer fail")
	}

	s = `"\`
	clean = keywords.CleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "" {
		panic("replacer fail")
	}

	s = "  a  b    c         de "
	clean = keywords.CleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "a b c de" {
		panic("replacer fail")
	}

	s = "\ta\nb \vc\fd\r\ne"
	clean = keywords.CleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "a b c d e" {
		panic("replacer fail")
	}

	s = "school music band club"
	clean = keywords.CleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "school music band club" {
		panic("breaking keywords")
	}
}
