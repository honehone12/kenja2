package endec

import "encoding/json"

type Json struct{}

func (j Json) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j Json) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j Json) ContentType() string {
	return "application/json; charset=utf-8"
}
