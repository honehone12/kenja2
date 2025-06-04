package documents

type Rating int32

const (
	RATING_UNKNOWN = iota
	RATING_ALL_AGES
	RATING_HENTAI
)

type ItemType int32

const (
	ITEM_TYPE_UNSPECIFIED = iota
	ITEM_TYPE_ANIME
	ITEM_TYPE_CHARACTER
)

type ItemId struct {
	Id       int64    `json:"id" bson:"id"`
	ItemType ItemType `json:"item_type" bson:"item_type"`
}

type Parent struct {
	Id           int64   `json:"id" bson:"id"`
	Name         string  `json:"name" bson:"name"`
	NameJapanese *string `json:"name_japanese,omitempty" bson:"name_japanese"`
}

type Candidate struct {
	Url          string   `json:"url" bson:"url"`
	Parent       *Parent  `json:"parent,omitempty" bson:"parent"`
	Name         string   `json:"name" bson:"name"`
	NameEnglish  *string  `json:"name_english,omitempty" bson:"name_english"`
	NameJapanese *string  `json:"name_japanese,omitempty" bson:"name_japanese"`
	Aliases      []string `json:"aliases,omitempty" bson:"aliases"`
}
