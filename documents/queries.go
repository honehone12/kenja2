package documents

type Rating int32

const (
	RATING_UNKNOWN Rating = iota
	RATING_ALL_AGES
	RATING_HENTAI
)

type TextQuery struct {
	Rating   Rating `json:"rating,omitempty" msgpack:"rating,omitempty"`
	Keywords string `json:"keywords,omitempty" msgpack:"keywords,omitempty"`
}

type VectorQuery struct {
}

type QueryResult struct {
	Candidates []Candidate `json:"result,omitempty" msgpack:"result,omitempty"`
}
