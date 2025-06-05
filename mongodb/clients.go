package mongodb

import (
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoClient struct {
	client   *mongo.Client
	database *mongo.Database
}

func connect(uri string) (mongoClient, error) {
	c := mongoClient{}

	db := os.Getenv("MONGO_DB_NAME")
	if len(db) == 0 {
		return c, errors.New("env for db name is not set")
	}

	ops := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ops)
	if err != nil {
		return c, err
	}

	database := client.Database(db)

	c.client = client
	c.database = database
	return c, nil
}

type collections struct {
	allAges *mongo.Collection
}

func newCollections(database *mongo.Database) (collections, error) {
	coll := collections{}

	all := os.Getenv("COLLECTION_ALL_AGES")
	if len(all) == 0 {
		return coll, errors.New("env for collection all ages is not set")
	}
	allAges := database.Collection(all)

	coll.allAges = allAges
	return coll, nil
}
