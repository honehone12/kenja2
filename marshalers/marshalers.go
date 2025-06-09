package marshalers

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Marshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	ContentType() string
}

func NewJsonMarshler() Json {
	return Json{}
}

func NewBsonMarshaler() Bson {
	return Bson{}
}

func NewMsgPackMarshaler() MsgPack {
	return MsgPack{}
}

type Json struct{}

func (j Json) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j Json) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j Json) ContentType() string {
	return "application/json; charset=utf8"
}

type Bson struct{}

func (b Bson) Marshal(v any) ([]byte, error) {
	return bson.Marshal(v)
}

func (b Bson) Unmarshal(data []byte, v any) error {
	return bson.Unmarshal(data, v)
}

func (b Bson) ContentType() string {
	return "application/bson"
}

type MsgPack struct{}

func (m MsgPack) Marshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m MsgPack) Unmarshal(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}

func (m MsgPack) ContentType() string {
	return "application/messagepack"
}
