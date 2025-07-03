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
		"id": "000000000000000000000000",
		"parent": {
			"id": "000000000000000000000001",
			"name": "Test Parent",
			"name_japanese": null
		},
		"name": "Test Data",
		"name_japanese": null
	}`

	c := Candidate{}
	if err := json.Unmarshal([]byte(s), &c); err != nil {
		panic(err)
	}

	oid := bson.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	if oid.Hex() != c.Parent.IdHex || c.Parent.Name != "Test Parent" || c.Parent.NameJapanese != "" {
		panic("decoded wrong parent")
	}
	if c.Name != "Test Data" || c.NameJapanese != "" {
		panic("decoded wrong names")
	}

	b, err := json.MarshalIndent(c, "\t", "\t")
	if err != nil {
		panic(err)
	}

	s = `{
		"id": "000000000000000000000000",
		"parent": {
			"id": "000000000000000000000001",
			"name": "Test Parent"
		},
		"name": "Test Data"
	}`
	if !bytes.Equal(b, []byte(s)) {
		fmt.Println(string(b))
		panic("encoded wrong candidate")
	}

	s = `{
		"id": "000000000000000000000000",
		"name": "Test Data"
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

func TestParentIsZero(t *testing.T) {
	p := Parent{}
	if !p.IsZero() {
		panic("parent should be zero")
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
		"item_type": 0,
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
