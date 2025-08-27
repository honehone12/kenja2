package documents

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type FlatDocument struct {
	Id           bson.ObjectID `json:"id" bson:"_id"`
	ItemType     ItemType      `json:"item_type" bson:"item_type"`
	Url          string        `json:"url,omitempty" bson:"url"`
	Unique       string        `json:"unique,omitempty" bson:"unique"`
	Img          string        `json:"img,omitempty" bson:"img"`
	Src          string        `json:"src,omitempty" bson:"src"`
	Parent       *Parent       `json:"parent,omitempty" bson:"parent"`
	Name         string        `json:"name,omitempty" bson:"name"`
	NameEnglish  string        `json:"name_english,omitempty" bson:"name_english"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese"`
}

type Parent struct {
	Name         string `json:"name" bson:"name"`
	NameJapanese string `json:"name_japanese,omitempty" bson:"name_japanese"`
	Url          string `json:"url,omitempty" bson:"url"`
}

func (p Parent) IsZero() bool {
	return len(p.Name) == 0 && len(p.NameJapanese) == 0 && len(p.Url) == 0
}

type SearchResult struct {
	Documents []FlatDocument `json:"documents,omitempty"`
}

type ImageVector struct {
	ImageVector bson.Vector `bson:"image_vector"`
}
