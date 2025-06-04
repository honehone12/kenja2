package mongodb

import (
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Collections struct {
	allAges *mongo.Collection
}

func NewCollections(database *mongo.Database) (Collections, error) {
	coll := Collections{}

	all := os.Getenv("COLLECTION_ALL_AGES")
	if len(all) == 0 {
		return coll, errors.New("env for collection all ages is not set")
	}
	allAges := database.Collection(all)

	coll.allAges = allAges
	return coll, nil
}
