package documents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestSeDe(t *testing.T) {
	s := `{
		"item_id": {
			"id": 0,
			"item_type": 1
		},
		"parent": {
			"id": 1,
			"name": "Test Parent",
			"name_japanese": null
		},
		"name": "Test Data",
		"name_english": "Test Name",
		"name_japanese": null
	}`

	c := Candidate{}
	if err := json.Unmarshal([]byte(s), &c); err != nil {
		panic(err)
	}
	if c.ItemId.Id != 0 || c.ItemId.ItemType != 1 {
		panic("decoded wrong item id")
	}
	if c.Parent.Id != 1 || c.Parent.Name != "Test Parent" || c.Parent.NameJapanese != nil {
		panic("decoded wrong parent")
	}
	if c.Name != "Test Data" || *c.NameEnglish != "Test Name" || c.NameJapanese != nil {
		panic("decoded wrong candidate")
	}

	b, err := json.MarshalIndent(c, "\t", "\t")
	if err != nil {
		panic(err)
	}

	s = `{
		"item_id": {
			"id": 0,
			"item_type": 1
		},
		"parent": {
			"id": 1,
			"name": "Test Parent"
		},
		"name": "Test Data",
		"name_english": "Test Name"
	}`
	if !bytes.Equal(b, []byte(s)) {
		fmt.Println(string(b))
		panic("encoded wrong candidate")
	}
}
