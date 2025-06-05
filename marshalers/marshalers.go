package marshalers

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Marshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

func NewJsonMarshler() Json {
	return Json{}
}

func NewBsonMarshaler() Bson {
	return Bson{}
}

type Json struct{}

func (j Json) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j Json) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

type Bson struct{}

func (b Bson) Marshal(v any) ([]byte, error) {
	return bson.Marshal(v)
}

func (b Bson) Unmarshal(data []byte, v any) error {
	return bson.Unmarshal(data, v)
}
