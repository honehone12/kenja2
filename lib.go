package kenja2

import (
	"errors"
	"kenja2/endec"
	"kenja2/engine"
	"kenja2/mongodb"
	"os"
)

func NewJson() endec.Json {
	return endec.Json{}
}

func NewMsgPack() endec.MsgPack {
	return endec.MsgPack{}
}

func ConnectAtlas[E endec.Encoder, D endec.Decoder](
	encoder E,
	decoder D,
) (engine.Engine, error) {
	uri := os.Getenv("SEARCHENGINE_URI")
	if len(uri) == 0 {
		return nil, errors.New("env for search engine uri is not set")
	}

	return mongodb.Connect(
		uri,
		encoder,
		decoder,
	)
}
