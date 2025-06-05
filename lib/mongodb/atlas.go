package mongodb

import (
	"context"
	"kenja2/lib/documents"
)

type AtlaSearch struct {
	mongoClient
	collections collections
}

func Connet(uri string) (*AtlaSearch, error) {
	m, err := connect(uri)
	if err != nil {
		return nil, err
	}

	coll, err := newCollections(m.database)
	if err != nil {
		return nil, err
	}

	return &AtlaSearch{
		mongoClient: m,
		collections: coll,
	}, nil
}

func (a *AtlaSearch) TextSearch(keywords string, rating documents.Rating) {

}

func (a *AtlaSearch) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
