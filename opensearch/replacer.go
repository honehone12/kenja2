package opensearch

import (
	"kenja2/keywords"
	"strings"
)

var __OPENSEARCH_REPLACER = strings.NewReplacer(
	"_", keywords.PLACE_HOLDER,
	"+", keywords.PLACE_HOLDER,
	"-", keywords.PLACE_HOLDER,
	"*", keywords.PLACE_HOLDER,
	"%", keywords.PLACE_HOLDER,
	"&", keywords.PLACE_HOLDER,
	"|", keywords.PLACE_HOLDER,
	"=", keywords.PLACE_HOLDER,
	"<", keywords.PLACE_HOLDER,
	">", keywords.PLACE_HOLDER,
	"!", keywords.PLACE_HOLDER,
	"?", keywords.PLACE_HOLDER,
	"^", keywords.PLACE_HOLDER,
	"~", keywords.PLACE_HOLDER,
	"{", keywords.PLACE_HOLDER,
	"}", keywords.PLACE_HOLDER,
	"[", keywords.PLACE_HOLDER,
	"]", keywords.PLACE_HOLDER,
	"(", keywords.PLACE_HOLDER,
	")", keywords.PLACE_HOLDER,
	":", keywords.PLACE_HOLDER,
	",", keywords.PLACE_HOLDER,
	".", keywords.PLACE_HOLDER,
	"/", keywords.PLACE_HOLDER,
	`\`, keywords.PLACE_HOLDER,
	`"`, keywords.PLACE_HOLDER,
	"'", keywords.PLACE_HOLDER,
	"`", keywords.PLACE_HOLDER,
)
