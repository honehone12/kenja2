package documents

import (
	"errors"
	"math/bits"
)

type Rating uint32

const (
	RATING_UNSPECIFIED Rating = iota
	RATING_ALL_AGES
	RATING_HENTAI
)

func (r Rating) Ok() error {
	if r > RATING_HENTAI {
		return errors.New("unexpected rating")
	}

	return nil
}

type VectorField uint32
type VFieldName string

const (
	VFIELD_UNSPECIFIED VectorField = iota
	VFIELD_TXT
	VFIELD_IMG
	VFIELD_STF
)

func (f VectorField) Name() (VFieldName, error) {
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

type ItemType uint32

const (
	ITEM_TYPE_UNSPECIFIED ItemType = iota
	ITEM_TYPE_ANIME
	ITEM_TYPE_CHARACTER
	ITEM_TYPE_ANIME_IMG_ONLY
	ITEM_TYPE_CHARACTER_IMG_ONLY
)

func (i ItemType) Ok() error {
	if i > ITEM_TYPE_CHARACTER_IMG_ONLY {
		return errors.New("unexpected item type")
	}

	return nil
}

func (i ItemType) IsAnime() bool {
	return i == ITEM_TYPE_ANIME || i == ITEM_TYPE_ANIME_IMG_ONLY
}

func (i ItemType) IsCharacter() bool {
	return i == ITEM_TYPE_CHARACTER || i == ITEM_TYPE_CHARACTER_IMG_ONLY
}

type FunctionId uint32

const (
	FID_UNSPECIFIED FunctionId = iota
	FID_PING
	FID_DOCUMENT
	FID_TEXTSEARCH_BUFFER
	FID_VECTORSEARCH_BUFFER
)

type FuzzyKind uint32

const (
	FUZZY_UNSPECIFIED = iota
	FUZZY_USE
)

func (f FuzzyKind) Ok() error {
	if f > FUZZY_USE {
		return errors.New("unexpected fuzzy kind")
	}

	return nil
}

type SearchQuery struct {
	FunctionId  FunctionId  `json:"function_id"`
	Rating      Rating      `json:"rating"`
	ItemType    ItemType    `json:"item_type"`
	VectorField VectorField `json:"vector_field"`
	Id          string      `json:"id"`
	Keywords    string      `json:"keywords"`
	FuzzyKind   FuzzyKind   `json:"fuzzy_kind"`
}

func (q *SearchQuery) StringSize() (uint64, bool) {
	len, carry := bits.Add64(uint64(len(q.Id)), uint64(len(q.Keywords)), 0)
	if carry > 0 {
		return 0, false
	}

	return len, true
}
