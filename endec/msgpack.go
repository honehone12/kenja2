package endec

import (
	"github.com/vmihailenco/msgpack/v5"
)

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
