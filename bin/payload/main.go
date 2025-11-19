package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"kenja2/documents"
)

func main() {
	q := documents.SearchQuery{
		FunctionId: 3,
		Keywords:   "ばっどがーる",
	}

	b, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}

	s := base64.RawURLEncoding.EncodeToString(b)
	fmt.Println(s)
}
