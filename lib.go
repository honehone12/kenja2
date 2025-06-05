package kenja2

import (
	"context"
	"encoding/json"
	"kenja2/marshalers"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Engine[M marshalers.Marshaler] interface {
	TextSearch(input []byte)
	Close(ctx context.Context) error
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
