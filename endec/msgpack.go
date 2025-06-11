package endec

import (
	"encoding/base64"

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

func (m MsgPack) String(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func (m MsgPack) FromString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
