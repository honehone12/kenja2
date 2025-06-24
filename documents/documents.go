package documents

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SingleDocument struct {
	Id           bson.ObjectID `json:"-" bson:"_id" msgpack:"-"`
	IdHex        string        `json:"id" bson:"-" msgpack:"id"`
	Parent       Parent        `json:"parent,omitzero" bson:"parent" msgpack:"parent,omitempty"`
	Name         string        `json:"name" bson:"name" msgpack:"name"`
	NameEnglish  string        `json:"name_english,omitempty" bson:"name_english" msgpack:"name_english,omitempty"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese" msgpack:"name_japanese,omitempty"`
	Url          string        `json:"url" bson:"url" msgpack:"url"`
	Studios      []string      `json:"studios,omitempty" bson:"studios" msgpack:"studios,omitempty"`
	Aliases      []string      `json:"aliases,omitempty" bson:"aliases" msgpack:"aliases,omitempty"`
}

// json encodes bson.ObjectId as string, but msgpack encodes as bytes
// this ensures ObjectId is encoded as string
func (d *SingleDocument) fillIdHex() {
	d.IdHex = d.Id.Hex()
	if !d.Parent.Id.IsZero() {
		d.Parent.IdHex = d.Parent.Id.Hex()
	}
}

func (d *SingleDocument) Format() {
	d.fillIdHex()
}

type Parent struct {
	Id           bson.ObjectID `json:"-" bson:"id" msgpack:"-"`
	IdHex        string        `json:"id" bson:"-" msgpack:"id"`
	Name         string        `json:"name" bson:"name" msgpack:"name"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese" msgpack:"name_japanese,omitempty"`
}

func (p Parent) IsZero() bool {
	return p.Id.IsZero() && len(p.Name) == 0 && len(p.NameJapanese) == 0
}

type Candidate struct {
	Id           bson.ObjectID `json:"-" bson:"_id" msgpack:"-"`
	IdHex        string        `json:"id" bson:"-" msgpack:"id"`
	Parent       Parent        `json:"parent,omitzero" bson:"parent" msgpack:"parent,omitempty"`
	Name         string        `json:"name" bson:"name" msgpack:"name"`
	NameJapanese string        `json:"name_japanese,omitempty" bson:"name_japanese" msgpack:"name_japanese,omitempty"`
}

// json encodes bson.ObjectId as string, but msgpack encodes as bytes
// this ensures ObjectId is encoded as string
func (c *Candidate) fillIdHex() {
	c.IdHex = c.Id.Hex()
	if !c.Parent.Id.IsZero() {
		c.Parent.IdHex = c.Parent.Id.Hex()
	}
}

func (c *Candidate) Format() {
	c.fillIdHex()
}

type SearchResult struct {
	Candidates []Candidate `json:"candidates,omitempty" msgpack:"candidates,omitempty"`
}

type Vector struct {
	TextVector  bson.Vector `bson:"text_vector"`
	ImageVector bson.Vector `bson:"image_vector"`
	StaffVector bson.Vector `bson:"staff_vector"`
}

func (v *Vector) BinaryField(field VectorField) (bson.Binary, error) {
	var b bson.Binary
	switch field {
	case VFIELD_TXT:
		b = v.TextVector.Binary()
	case VFIELD_IMG:
		b = v.ImageVector.Binary()
	case VFIELD_STF:
		b = v.StaffVector.Binary()
	default:
		return b, errors.New("unexpected field")
	}

	if b.IsZero() {
		return b, errors.New("specified field is empty")
	}

	return b, nil
}
