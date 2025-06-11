package endec

type Encoder interface {
	Marshal(v any) ([]byte, error)
	String(b []byte) string
	ContentType() string
}

type Decoder interface {
	Unmarshal(data []byte, v any) error
	FromString(s string) ([]byte, error)
	ContentType() string
}
