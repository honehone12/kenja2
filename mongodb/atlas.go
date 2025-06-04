package mongodb

import (
	"context"
	"errors"
	"kenja2/documents"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AtlaSearch struct {
	client      *mongo.Client
	database    *mongo.Database
	collections Collections
}

func Connet(uri string) (*AtlaSearch, error) {
	db := os.Getenv("MONGO_DB_NAME")
	if len(db) == 0 {
		return nil, errors.New("env for db name is not set")
	}

	ops := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ops)
	if err != nil {
		return nil, err
	}

	database := client.Database(db)
	collections, err := NewCollections(database)
	if err != nil {
		return nil, err
	}

	return &AtlaSearch{
		client,
		database,
		collections,
	}, nil
}

func (a *AtlaSearch) TextSearch(keywords string, rating documents.Rating) {

}

func (a *AtlaSearch) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
