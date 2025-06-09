package documents

type Rating int32

const (
	RATING_UNKNOWN Rating = iota
	RATING_ALL_AGES
	RATING_HENTAI
)

type TextQuery struct {
	Rating   Rating `json:"rating,omitempty" bson:"rating,omitempty" msgpack:"rating,omitempty"`
	Keywords string `json:"keywords,omitempty" bson:"keywords,omitempty" msgpack:"keywords,omitempty"`
}

type QueryResult struct {
	Candidates []Candidate `json:"result,omitempty" bson:"result,omitempty" msgpack:"result,omitempty"`
}
