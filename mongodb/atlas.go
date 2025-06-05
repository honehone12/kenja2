package mongodb

import (
	"context"
	"kenja2/marshalers"
)

type Atlas[E, D marshalers.Marshaler] struct {
	mongoClient
	collections collections
	encoder     E
	decoder     D
}

func Connet[E, D marshalers.Marshaler](
	uri string,
	encoder E,
	decoder D,
) (*Atlas[E, D], error) {
	mongoClient, err := connect(uri)
	if err != nil {
		return nil, err
	}

	collections, err := newCollections(mongoClient.database)
	if err != nil {
		return nil, err
	}

	return &Atlas[E, D]{
		mongoClient,
		collections,
		encoder,
		decoder,
	}, nil
}

func (a *Atlas[E, D]) TextSearch(input []byte) ([]byte, error) {
	return nil, nil
}

func (a *Atlas[E, D]) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
