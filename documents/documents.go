package documents

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FlatDocument struct {
	Id           bson.ObjectID `json:"id" bson:"_id"`
	ItemType     ItemType      `json:"item_type" bson:"item_type"`
	Url          string        `json:"url,omitempty" bson:"url"`
	Img          string        `json:"img,omitempty" bson:"img"`
	License      []string      `json:"license,omitempty" bson:"license"`
	Src          string        `json:"src,omitempty" bson:"src"`
	Parent       *Parent       `json:"parent,omitempty" bson:"parent"`
	Name         string        `json:"name,omitempty" bson:"name"`
	NameEnglish  string        `json:"name_english,omitempty" bson:"name_english"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese"`
}

type Parent struct {
	Id           bson.ObjectID `json:"id" bson:"id"`
	Name         string        `json:"name" bson:"name"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese"`
}

func (p Parent) IsZero() bool {
	return p.Id.IsZero() && len(p.Name) == 0 && len(p.NameJapanese) == 0
}

type SearchResult struct {
	Documents []FlatDocument `json:"documents,omitempty"`
}

type ImageVector struct {
	ImageVector bson.Vector `bson:"image_vector"`
}

func (v *ImageVector) Binary() (bson.Binary, error) {
	b := v.ImageVector.Binary()
	if b.IsZero() {
		return b, errors.New("iamge vector fiel is empty")
	}
	return b, nil
}
