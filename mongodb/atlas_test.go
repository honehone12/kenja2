package mongodb

import (
	"kenja2/documents"
	"kenja2/engine"
)

var _ engine.Engine = &Atlas{}

var _ engine.Stream[documents.FlatDocument] = &Stream{}
