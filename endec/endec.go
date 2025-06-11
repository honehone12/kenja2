package endec

type Encoder interface {
	Marshal(v any) ([]byte, error)
	ContentType() string
}

type Decoder interface {
	Unmarshal(data []byte, v any) error
	ContentType() string
}
