package mongodb

import (
	"context"
	"kenja2/documents"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Stream struct {
	stream *mongo.Cursor
}

func (s *Stream) Next(ctx context.Context) bool {
	return s.stream.Next(ctx)
}

func (s *Stream) Decode(val *documents.FlatDocument) error {
	return s.stream.Decode(val)
}

func (s *Stream) Err() error {
	return s.stream.Err()
}

func (s *Stream) Close(ctx context.Context) error {
	return s.stream.Close(ctx)
}
