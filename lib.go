package kenja2

import (
	"context"
	"errors"
	"fmt"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/mongodb"
	"os"
)

func ConnectAtlas(maxConn uint64, logger engine.Logger) (engine.Engine, error) {
	uri := os.Getenv("SEARCHENGINE_URI")
	if len(uri) == 0 {
		return nil, errors.New("env for search engine uri is not set")
	}

	return mongodb.Connect(uri, maxConn, logger)
}

func ConnectAtlasAws(maxConn uint64, logger engine.Logger) (engine.Engine, error) {
	uri := os.Getenv("SEARCHENGINE_URI")
	if len(uri) == 0 {
		return nil, errors.New("env for search engine uri is not set")
	}

	return mongodb.ConnectAws(uri, maxConn, logger)
}

func BufferedQuery(
	engine engine.Engine,
	ctx context.Context,
	q *documents.SearchQuery,
) ([]documents.FlatDocument, error) {
	switch q.FunctionId {
	case documents.FID_PING:
		return nil, fmt.Errorf("function: %d is not in service", q.FunctionId)
	case documents.FID_RANDOM:
		return engine.Random(ctx, q)
	case documents.FID_TEXTSEARCH:
		return engine.TextSearch(ctx, q)
	case documents.FID_VECTORSEARCH:
		return engine.VectorSeach(ctx, q)
	case documents.FID_DIRECT_VECTORSEARCH:
		return engine.DirectVectorSearch(ctx, q)
	case documents.FID_MULTI_VECTOR_SEARCH:
		return engine.MultiVectorSearch(ctx, q)
	default:
		return nil, errors.New("invalid function id")
	}
}

func StreamQuery(
	engine engine.Engine,
	ctx context.Context,
	q *documents.SearchQuery,
) (engine.Stream[documents.FlatDocument], error) {
	switch q.FunctionId {
	case documents.FID_TEXTSEARCH_STREAM:
		return engine.TextSearchStream(ctx, q)
	case documents.FID_VECTORSEARCH_STREAM:
		return engine.VectorSearchStream(ctx, q)
	case documents.FID_DIRECT_VECTORSEARCH_STREAM:
		return engine.DirectVectorSearchStream(ctx, q)
	case documents.FID_MULTI_VECTOR_SEARCH_STREAM:
		return engine.MultiVectorSearchStream(ctx, q)
	default:
		return nil, errors.New("invalid function id")
	}
}
