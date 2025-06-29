package mongodb

import (
	"kenja2/keywords"
	"strings"
)

var __MONGO_REPLACER = strings.NewReplacer(
	"$", keywords.PLACE_HOLDER,
	".", keywords.PLACE_HOLDER,
	",", keywords.PLACE_HOLDER,
	"{", keywords.PLACE_HOLDER,
	"}", keywords.PLACE_HOLDER,
	"[", keywords.PLACE_HOLDER,
	"]", keywords.PLACE_HOLDER,
	"(", keywords.PLACE_HOLDER,
	")", keywords.PLACE_HOLDER,
	":", keywords.PLACE_HOLDER,
	";", keywords.PLACE_HOLDER,
	"/", keywords.PLACE_HOLDER,
	"`", keywords.PLACE_HOLDER,
	"'", keywords.PLACE_HOLDER,
	`"`, keywords.PLACE_HOLDER,
)
