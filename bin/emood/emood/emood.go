package emood

import "go.mongodb.org/mongo-driver/v2/bson"

type Emood struct {
	Src    string
	Tensor bson.Vector
}
