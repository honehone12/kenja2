package documents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/vmihailenco/msgpack/v5"
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
	if c.ItemType != 1 {
		panic("decoded wrong item type")
	}
	if c.Url != "test.kenja" {
		panic("decoded wrong url")
	}
	oid := bson.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	if !bytes.Equal(oid[:], c.Parent.Id[:]) || c.Parent.Name != "Test Parent" || c.Parent.NameJapanese != "" {
		panic("decoded wrong parent")
	}
	if c.Name != "Test Data" || c.NameEnglish != "Test Name" || c.NameJapanese != "" {
		panic("decoded wrong names")
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

func TestParentIsZero(t *testing.T) {
	p := Parent{}
	if !p.IsZero() {
		panic("parent should be zero")
	}
}

func TestDocumentsMsgPack(t *testing.T) {
	id := bson.NewObjectID()
	c := Candidate{
		ItemType: 1,
		Url:      "kenja.test",
		Parent: Parent{
			Id:           id,
			Name:         "Test",
			NameJapanese: "",
		},
		Name:         "Test",
		NameEnglish:  "Test",
		NameJapanese: "",
		Aliases:      []string{"Test"},
	}
	b, err := msgpack.Marshal(c)
	if err != nil {
		panic(err)
	}

	c = Candidate{}
	if err := msgpack.Unmarshal(b, &c); err != nil {
		panic(err)
	}
	if c.ItemType != 1 {
		panic("decoded wrong item type")
	}
	if c.Url != "kenja.test" {
		panic("decoded wrong url")
	}
	if !bytes.Equal(c.Parent.Id[:], id[:]) || c.Parent.Name != "Test" || c.Parent.NameJapanese != "" {
		panic("decoded wrong parent")
	}
	if c.Name != "Test" || c.NameEnglish != "Test" || c.NameJapanese != "" {
		panic("decoded wrong names")
	}
	if len(c.Aliases) != 1 || c.Aliases[0] != "Test" {
		panic("decoded wrong aliases")
	}

	c.Parent = Parent{}
	c.Aliases = nil
	b, err = msgpack.Marshal(c)
	if err != nil {
		panic(err)
	}

	c = Candidate{}
	if err = msgpack.Unmarshal(b, &c); err != nil {
		panic(err)
	}
	if !bytes.Equal(c.Parent.Id[:], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) || c.Parent.Name != "" || c.Parent.NameJapanese != "" {
		panic("decoded wrong parent")
	}
	if c.Aliases != nil {
		panic("decoded wrong aliases")
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

func TestQueriesMsgPack(t *testing.T) {
	q := TextQuery{
		Rating:   1,
		Keywords: "miku miku",
	}
	b, err := msgpack.Marshal(q)
	if err != nil {
		panic(err)
	}

	q = TextQuery{}
	if err = msgpack.Unmarshal(b, &q); err != nil {
		panic(err)
	}
	if q.Rating != 1 || q.Keywords != "miku miku" {
		panic("decoded wrong rating")
	}
}

func TestJsonBsonMsgPack(t *testing.T) {
	d := make([]Candidate, 10000)
	r := QueryResult{Candidates: d}

	now1 := time.Now()

	jb, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	rj := QueryResult{}
	if err = json.Unmarshal(jb, &rj); err != nil {
		panic(err)
	}

	now2 := time.Now()
	fmt.Printf("json: %d bytes\n", len(jb))
	fmt.Printf("time: %dmilsecs\n", now2.Sub(now1).Milliseconds())

	bb, err := bson.Marshal(r)
	if err != nil {
		panic(err)
	}
	rb := QueryResult{}
	if err = bson.Unmarshal(bb, &rb); err != nil {
		panic(err)
	}

	now3 := time.Now()
	fmt.Printf("bson: %d bytes\n", len(bb))
	fmt.Printf("time: %dmilsecs\n", now3.Sub(now2).Milliseconds())

	mb, err := msgpack.Marshal(r)
	if err != nil {
		panic(err)
	}
	rm := QueryResult{}
	if err = msgpack.Unmarshal(mb, &rm); err != nil {
		panic(err)
	}

	now4 := time.Now()
	fmt.Printf("msgpack: %d bytes\n", len(mb))
	fmt.Printf("time: %dmilsecs\n", now4.Sub(now3).Milliseconds())
}
