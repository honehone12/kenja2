package documents

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Parent struct {
	Id           bson.ObjectID `json:"id" bson:"id" msgpack:"id"`
	Name         string        `json:"name" bson:"name" msgpack:"name"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese" msgpack:"name_japanese,omitempty"`
}

func (p Parent) IsZero() bool {
	return p.Id.IsZero() && len(p.Name) == 0 && len(p.NameJapanese) == 0
}

type Candidate struct {
	Id           bson.ObjectID `json:"id" bson:"_id" msgpack:"id"`
	Url          string        `json:"url" bson:"url" msgpack:"url"`
	Parent       Parent        `json:"parent,omitzero" bson:"parent" msgpack:"parent,omitempty"`
	Name         string        `json:"name" bson:"name" msgpack:"name"`
	NameEnglish  string        `json:"name_english,omitempty" bson:"name_english" msgpack:"name_english,omitempty"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese" msgpack:"name_japanese,omitempty"`
	Aliases      []string      `json:"aliases,omitempty" bson:"aliases" msgpack:"aliases,omitempty"`
}

type Vector struct {
	TextVector  bson.Vector `bson:"text_vector"`
	ImageVector bson.Vector `bson:"image_vector"`
}

func (v *Vector) BinaryField(field VectorField) (bson.Binary, error) {
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
