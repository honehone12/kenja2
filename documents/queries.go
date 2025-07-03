package documents

import "errors"

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

type Query interface {
	Size() int
}

type DocumentQuery struct {
	Id string `json:"id" msgpack:"id"`
}

func (d *DocumentQuery) Size() int {
	if d == nil {
		return 0
	}

	return len(d.Id)
}

type TextQuery struct {
	Rating   Rating   `json:"rating" msgpack:"rating"`
	ItemType ItemType `json:"item_type" msgpack:"item_type"`
	Keywords string   `json:"keywords" msgpack:"keywords"`
}

func (t *TextQuery) Size() int {
	if t == nil {
		return 0
	}

	return len(t.Keywords) + 2
}

type VectorQuery struct {
	Rating      Rating      `json:"rating" msgpack:"rating"`
	ItemType    ItemType    `json:"item_type" msgpack:"item_type"`
	SourceField VectorField `json:"source_field" msgpack:"source_field"`
	TargetField VectorField `json:"target_field" msgpack:"target_field"`
	Id          string      `json:"id" msgpack:"id"`
}

func (v *VectorQuery) Size() int {
	if v == nil {
		return 0
	}

	return len(v.Id) + 4
}
