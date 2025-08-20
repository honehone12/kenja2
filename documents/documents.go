package documents

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FlatDocument struct {
	Id           bson.ObjectID `json:"id" bson:"_id"`
	ItemType     ItemType      `json:"item_type" bson:"item_type"`
	Rating       Rating        `json:"rating" bson:"rating"`
	Url          string        `json:"url" bson:"url"`
	Img          string        `json:"img" bson:"img"`
	Parent       Parent        `json:"parent,omitzero" bson:"parent"`
	Name         string        `json:"name" bson:"name"`
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

type Vector struct {
	ImageVector bson.Vector `bson:"image_vector"`
}

func (v *Vector) Binary(field VectorField) (bson.Binary, error) {
	var b bson.Binary
	switch field {
	case VFIELD_IMG:
		b = v.ImageVector.Binary()
	// add fields here when use
	default:
		return b, errors.New("unexpected field")
	}

	if b.IsZero() {
		return b, errors.New("specified field is empty")
	}

	return b, nil
}
