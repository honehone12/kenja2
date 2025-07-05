package documents

import (
	"errors"
	"math/bits"
)

type Rating int32

const (
	RATING_UNSPECIFIED Rating = iota
	RATING_ALL_AGES
	RATING_HENTAI
)
const RATING_MAX = RATING_HENTAI

func (r Rating) Ok() error {
	if r < 0 || r > RATING_MAX {
		return errors.New("unexpected rating")
	}

	return nil
}

type VectorField int32

const (
	VFIELD_UNSPECIFIED VectorField = iota
	VFIELD_TXT
	VFIELD_IMG
	VFIELD_STF
)

func (f VectorField) String() (string, error) {
	switch f {
	case VFIELD_TXT:
		return "text_vector", nil
	case VFIELD_IMG:
		return "image_vector", nil
	case VFIELD_STF:
		return "staff_vector", nil
	default:
		return "", errors.New("unexpected vector fields option")
	}
}

type ItemType int32

const (
	ITEM_TYPE_UNSPECIFIED ItemType = iota
	ITEM_TYPE_ANIME
	ITEM_TYPE_CHARACTER
)
const ITEM_TYPE_MAX = ITEM_TYPE_CHARACTER

func (i ItemType) Ok() error {
	if i < 0 || i > ITEM_TYPE_MAX {
		return errors.New("unexpected item type")
	}

	return nil
}

type SearchQuery struct {
	Rating      Rating      `json:"rating" msgpack:"rating"`
	ItemType    ItemType    `json:"item_type" msgpack:"item_type"`
	SourceField VectorField `json:"source_field" msgpack:"source_field"`
	TargetField VectorField `json:"target_field" msgpack:"target_field"`
	Id          string      `json:"id" msgpack:"id"`
	Keywords    string      `json:"keywords" msgpack:"keywords"`
}

func (q *SearchQuery) StringSize() (uint64, bool) {
	len, carry := bits.Add64(uint64(len(q.Id)), uint64(len(q.Keywords)), 0)
	if carry > 0 {
		return 0, false
	}

	return len, true
}
