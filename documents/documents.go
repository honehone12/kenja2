package documents

import "go.mongodb.org/mongo-driver/v2/bson"

type ItemType int32

const (
	ITEM_TYPE_UNSPECIFIED ItemType = iota
	ITEM_TYPE_ANIME
	ITEM_TYPE_CHARACTER
)

type Parent struct {
	Id           bson.ObjectID `json:"id" bson:"id" msgpack:"id"`
	Name         string        `json:"name" bson:"name" msgpack:"name"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese,omitempty" msgpack:"name_japanese,omitempty"`
}

func (p Parent) IsZero() bool {
	return p.Id.IsZero() && len(p.Name) == 0 && len(p.NameJapanese) == 0
}

type Candidate struct {
	ItemType     ItemType `json:"item_type,omitempty" bson:"item_type,omitempty" msgpack:"item_type,omitempty"`
	Url          string   `json:"url" bson:"url" msgpack:"url"`
	Parent       Parent   `json:"parent,omitzero" bson:"parent,omitempty" msgpack:"parent,omitempty"`
	Name         string   `json:"name" bson:"name" msgpack:"name"`
	NameEnglish  string   `json:"name_english,omitempty" bson:"name_english,omitempty" msgpack:"name_english,omitempty"`
	NameJapanese string   `json:"name_japanese,omitempty" bson:"name_japanese,omitempty" msgpack:"name_japanese,omitempty"`
	Aliases      []string `json:"aliases,omitempty" bson:"aliases,omitempty" msgpack:"aliases,omitempty"`
}
