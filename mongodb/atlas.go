package mongodb

import (
	"context"
	"kenja2/documents"
	"kenja2/marshalers"
)

type Atlas struct {
	mongoClient
	collections collections
	marshaler   marshalers.Marshaler
}

func Connet[M marshalers.Marshaler](uri string, marshaler M) (*Atlas, error) {
	mongoClient, err := connect(uri)
	if err != nil {
		return nil, err
	}

	collections, err := newCollections(mongoClient.database)
	if err != nil {
		return nil, err
	}

	return &Atlas{
		mongoClient,
		collections,
		marshaler,
	}, nil
}

func (a *Atlas) TextSearch(keywords string, rating documents.Rating) {

}

func (a *Atlas) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
