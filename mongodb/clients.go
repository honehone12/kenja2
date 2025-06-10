package mongodb

import (
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoClient struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
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

	coll := os.Getenv("MONGO_COLLECTION")
	if len(coll) == 0 {
		return c, errors.New("env for collection is not set")
	}

	collection := database.Collection(coll)

	c.client = client
	c.database = database
	c.collection = collection
	return c, nil
}
