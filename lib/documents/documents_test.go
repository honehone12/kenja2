package documents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestDocuments(t *testing.T) {
	s := `{
		"url": "test.kenja",
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
	if c.Url != "test.kenja" {
		panic("decoded wrong url")
	}
	if c.Parent.Id != 1 || c.Parent.Name != "Test Parent" || c.Parent.NameJapanese != nil {
		panic("decoded wrong parent")
	}
	if c.Name != "Test Data" || *c.NameEnglish != "Test Name" || c.NameJapanese != nil {
		panic("decoded wrong candidate")
	}
	if c.Aliases != nil {
		panic("decoded wrong aliases")
	}

	b, err := json.MarshalIndent(c, "\t", "\t")
	if err != nil {
		panic(err)
	}

	s = `{
		"url": "test.kenja",
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
