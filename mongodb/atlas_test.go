package mongodb

import (
	"kenja2/engine"
)

var _ engine.BufferingEngine = &Atlas{}
