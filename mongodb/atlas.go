package mongodb

import (
	"context"
	"kenja2/documents"
	"kenja2/marshalers"
)

type Atlas struct {
	mongoClient
	collections collections
	encoder     marshalers.Marshaler
	decoder     marshalers.Marshaler
}

func Connet[E, D marshalers.Marshaler](
	uri string,
	encoder E,
	decoder D,
) (*Atlas, error) {
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
		encoder,
		decoder,
	}, nil
}

func (a *Atlas) TextSearch(keywords string, rating documents.Rating) {

}

func (a *Atlas) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
