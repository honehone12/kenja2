package documents

import (
	"fmt"
	"math/bits"
)

type ItemType uint32

const (
	ITEM_TYPE_UNSPECIFIED ItemType = iota
	ITEM_TYPE_ANIME
	ITEM_TYPE_CHARACTER
	ITEM_TYPE_YVIDEO
)

func (i ItemType) Ok() error {
	if i > ITEM_TYPE_YVIDEO {
		return fmt.Errorf("unexpected item type: %d", i)
	}
	return nil
}

type FunctionId uint32

const (
	FID_UNSPECIFIED FunctionId = iota
	FID_PING
	FID_NEW_ITEMS
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
		return fmt.Errorf("unexpected fuzzy kind: %d", f)
	}
	return nil
}

type SearchQuery struct {
	FunctionId FunctionId `json:"function_id"`
	ItemType   ItemType   `json:"item_type"`
	Id         string     `json:"id"`
	Keywords   string     `json:"keywords"`
	FuzzyKind  FuzzyKind  `json:"fuzzy_kind"`
}

func (q *SearchQuery) StringSize() (uint64, bool) {
	len, carry := bits.Add64(uint64(len(q.Id)), uint64(len(q.Keywords)), 0)
	if carry > 0 {
		return 0, false
	}

	return len, true
}
