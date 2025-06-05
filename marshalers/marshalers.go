package marshalers

type Marshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
