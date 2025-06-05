package mongodb

import (
	"context"
	"kenja2/documents"
	"kenja2/marshalers"
)

type AtlaSearch struct {
	mongoClient
	collections collections
	marshaler   marshalers.Marshaler
}

func Connet[M marshalers.Marshaler](uri string, marshaler M) (*AtlaSearch, error) {
	mongoClient, err := connect(uri)
	if err != nil {
		return nil, err
	}

	collections, err := newCollections(mongoClient.database)
	if err != nil {
		return nil, err
	}

	return &AtlaSearch{
		mongoClient,
		collections,
		marshaler,
	}, nil
}

func (a *AtlaSearch) TextSearch(keywords string, rating documents.Rating) {

}

func (a *AtlaSearch) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
