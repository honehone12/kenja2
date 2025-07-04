package documents

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var _ bson.Zeroer = Parent{}

func TestParentIsZero(t *testing.T) {
	p := Parent{}
	if !p.IsZero() {
		panic("parent should be zero")
	}
}
