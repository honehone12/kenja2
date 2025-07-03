package kenja2

import (
	"errors"
	"kenja2/engine"
	"kenja2/mongodb"
	"os"
)

func ConnectAtlas(logger engine.Logger) (engine.Engine, error) {
	uri := os.Getenv("SEARCHENGINE_URI")
	if len(uri) == 0 {
		return nil, errors.New("env for search engine uri is not set")
	}

	return mongodb.Connect(uri, logger)
}

func ConnectAtlasAws(logger engine.Logger) (engine.Engine, error) {
	uri := os.Getenv("SEARCHENGINE_URI")
	if len(uri) == 0 {
		return nil, errors.New("env for search engine uri is not set")
	}

	return mongodb.ConnectAws(uri, logger)
}
