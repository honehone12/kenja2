package documents

type Rating int32

const (
	RATING_UNKNOWN Rating = iota
	RATING_ALL_AGES
	RATING_HENTAI
)

type TextQuery struct {
	Rating   Rating `json:"rating" bson:"rating"`
	Keywords string `json:"keywords" bson:"keywords"`
}
