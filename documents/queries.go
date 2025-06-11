package documents

import "errors"

type Rating byte

const (
	RATING_UNSPECIFIED Rating = iota
	RATING_ALL_AGES
	RATING_HENTAI
)
const RATING_MAX = int32(RATING_HENTAI)

func (r Rating) I32() (int32, error) {
	i32 := int32(r)
	if i32 < 0 || i32 > RATING_MAX {
		return 0, errors.New("unexpected rating")
	}

	return i32, nil
}

type VectorField byte

const (
	VFIELD_UNSPECIFIED VectorField = iota
	VFIELD_TXT
	VFIELD_IMG
)

func (f VectorField) String() (string, error) {
	switch f {
	case VFIELD_TXT:
		return "text_vector", nil
	case VFIELD_IMG:
		return "image_vector", nil
	default:
		return "", errors.New("unexpected vector fields option")
	}
}

type ItemType byte

const (
	ITEM_TYPE_UNSPECIFIED ItemType = iota
	ITEM_TYPE_ANIME
	ITEM_TYPE_CHARACTER
)
const ITEM_TYPE_MAX = int32(ITEM_TYPE_CHARACTER)

func (i ItemType) I32() (int32, error) {
	i32 := int32(i)
	if i32 < 0 || i32 > ITEM_TYPE_MAX {
		return 0, errors.New("unexpected item type")
	}

	return i32, nil
}

type TextQuery struct {
	Rating   Rating   `json:"rating" msgpack:"rating"`
	ItemType ItemType `json:"item_type" msgpack:"item_type"`
	Keywords string   `json:"keywords" msgpack:"keywords"`
}

type VectorQuery struct {
	Rating      Rating      `json:"rating" msgpack:"rating"`
	ItemType    ItemType    `json:"item_type" msgpack:"item_type"`
	SourceField VectorField `json:"source_field" msgpack:"source_field"`
	TargetField VectorField `json:"target_field" msgpack:"target_field"`
	Id          string      `json:"id" msgpack:"id"`
}

type QueryResult struct {
	Candidates []Candidate `json:"candidates,omitempty" msgpack:"candidates,omitempty"`
}
