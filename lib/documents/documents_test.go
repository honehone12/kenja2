package documents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var _ bson.Zeroer = Parent{}

func TestDocuments(t *testing.T) {
	s := `{
		"item_type": 1,
		"url": "test.kenja",
		"parent": {
			"id": "000000000000000000000001",
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
	if c.ItemType != 1 {
		panic("decoded wrong item type")
	}
	oid := bson.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	if !bytes.Equal(oid[:], c.Parent.Id[:]) || c.Parent.Name != "Test Parent" || c.Parent.NameJapanese != "" {
		panic("decoded wrong parent")
	}
	if c.Name != "Test Data" || c.NameEnglish != "Test Name" || c.NameJapanese != "" {
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
		"item_type": 1,
		"url": "test.kenja",
		"parent": {
			"id": "000000000000000000000001",
			"name": "Test Parent"
		},
		"name": "Test Data",
		"name_english": "Test Name"
	}`
	if !bytes.Equal(b, []byte(s)) {
		fmt.Println(string(b))
		panic("encoded wrong candidate")
	}

	s = `{
		"item_type": 1,
		"url": "test.kenja",
		"name": "Test Data",
		"name_english": "Test Name"
	}`
	c.Parent = Parent{}
	b, err = json.MarshalIndent(c, "\t", "\t")
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(b, []byte(s)) {
		fmt.Println(string(b))
		panic("encoded wrong candidate")
	}
}

func TestQueries(t *testing.T) {
	s := `{
		"rating": 1,
		"keywords": "miku miku"
	}`

	q := TextQuery{}
	if err := json.Unmarshal([]byte(s), &q); err != nil {
		panic(err)
	}
	if q.Rating != 1 || q.Keywords != "miku miku" {
		panic("decoded wrong query")
	}

	q.Rating = 2
	q.Keywords = "muku muku"
	s = `{
		"rating": 2,
		"keywords": "muku muku"
	}`
	b, err := json.MarshalIndent(q, "\t", "\t")
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(b, []byte(s)) {
		fmt.Println(string(b))
		panic("encoded wrong query")
	}
}
