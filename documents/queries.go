package documents

import "errors"

type Rating byte

const (
	RATING_UNKNOWN Rating = iota
	RATING_ALL_AGES
	RATING_HENTAI
)

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

type ItemTypeQuery byte

func (i ItemTypeQuery) To32() (ItemType, error) {
	i32 := ItemType(i)
	if i32 < 0 || i32 > ITEM_TYPE_MAX {
		return 0, errors.New("unexpected item type query")
	}

	return i32, nil
}

type TextQuery struct {
	Rating   Rating        `json:"rating,omitempty" msgpack:"rating,omitempty"`
	ItemType ItemTypeQuery `json:"item_type,omitempty" msgpack:"item_type,omitempty"`
	Keywords string        `json:"keywords,omitempty" msgpack:"keywords,omitempty"`
}

type VectorQuery struct {
	Rating      Rating        `json:"rating,omitempty" msgpack:"rating,omitempty"`
	ItemType    ItemTypeQuery `json:"item_type,omitempty" msgpack:"item_type,omitempty"`
	SourceField VectorField   `json:"source_field,omitempty" msgpack:"source_field,omitempty"`
	TargetField VectorField   `json:"target_field,omitempty" msgpack:"target_field,omitempty"`
	Id          string        `json:"id,omitempty" msgpack:"id,omitempty"`
}

type QueryResult struct {
	Candidates []Candidate `json:"result,omitempty" msgpack:"result,omitempty"`
}
