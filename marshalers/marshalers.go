package marshalers

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

type Marshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	ContentType() string
}

func NewJsonMarshler() Json {
	return Json{}
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
