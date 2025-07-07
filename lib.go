package kenja2

import (
	"context"
	"errors"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/mongodb"
	"os"
)

func ConnectAtlas(maxConn uint64, logger engine.Logger) (engine.BufferingEngine, error) {
	uri := os.Getenv("SEARCHENGINE_URI")
	if len(uri) == 0 {
		return nil, errors.New("env for search engine uri is not set")
	}

	return mongodb.Connect(uri, maxConn, logger)
}

func ConnectAtlasAws(maxConn uint64, logger engine.Logger) (engine.BufferingEngine, error) {
	uri := os.Getenv("SEARCHENGINE_URI")
	if len(uri) == 0 {
		return nil, errors.New("env for search engine uri is not set")
	}

	return mongodb.ConnectAws(uri, maxConn, logger)
}

func BufferedQuery(
	engine engine.BufferingEngine,
	ctx context.Context,
	q *documents.SearchQuery,
) ([]documents.FlatDocument, error) {
	switch q.FunctionId {
	case documents.PING:
		return nil, engine.Ping(ctx)
	case documents.DOCUMENT:
		return engine.Document(ctx, q)
	case documents.TEXTSEARCH_BUFFER:
		return engine.TextSearch(ctx, q)
	case documents.VECTORSEARCH_BUFFER:
		return engine.VectorSeach(ctx, q)
	default:
		return nil, errors.New("invalid function id")
	}
}
