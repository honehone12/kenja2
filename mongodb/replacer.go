package mongodb

import (
	"kenja2/keywords"
	"strings"
)

var __mongo_replacer = strings.NewReplacer(
	"/", keywords.PLACE_HOLDER_SPACE,
	"$", keywords.PLACE_HOLDER_NONE,
	".", keywords.PLACE_HOLDER_NONE,
	",", keywords.PLACE_HOLDER_NONE,
	"{", keywords.PLACE_HOLDER_NONE,
	"}", keywords.PLACE_HOLDER_NONE,
	"[", keywords.PLACE_HOLDER_NONE,
	"]", keywords.PLACE_HOLDER_NONE,
	"(", keywords.PLACE_HOLDER_NONE,
	")", keywords.PLACE_HOLDER_NONE,
	":", keywords.PLACE_HOLDER_NONE,
	";", keywords.PLACE_HOLDER_NONE,
	"`", keywords.PLACE_HOLDER_NONE,
	"'", keywords.PLACE_HOLDER_NONE,
	`"`, keywords.PLACE_HOLDER_NONE,
	`\`, keywords.PLACE_HOLDER_NONE,
)
