package documents

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FlatDocument struct {
	Id           bson.ObjectID `json:"id" bson:"_id"`
	ItemType     ItemType      `json:"item_type" bson:"item_type"`
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
	TextVector  bson.Vector `bson:"text_vector"`
	ImageVector bson.Vector `bson:"image_vector"`
}

func (v *Vector) Binary(field VectorField) (bson.Binary, error) {
	var b bson.Binary
	switch field {
	case VFIELD_TXT:
		b = v.TextVector.Binary()
	case VFIELD_IMG:
		b = v.ImageVector.Binary()
	default:
		return b, errors.New("unexpected field")
	}

	if b.IsZero() {
		return b, errors.New("specified field is empty")
	}

	return b, nil
}

func (v *Vector) Float32Array(field VectorField) ([]float32, error) {
	var f []float32
	var ok bool
	switch field {
	case VFIELD_TXT:
		f, ok = v.TextVector.Float32OK()
	case VFIELD_IMG:
		f, ok = v.ImageVector.Float32OK()
	default:
		return nil, errors.New("unexpected field")
	}

	if len(f) == 0 {
		return nil, errors.New("specified field is empty")
	}
	if !ok {
		return nil, errors.New("failed to convert to float32 array")
	}

	return f, nil
}
