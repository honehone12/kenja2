package lambda

import (
	"context"
	"fmt"
	"kenja2"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/lambda/logs"
)

type HadlerFunc[
	Q documents.Query,
	R any,
] = func(context.Context, Q) (R, error)
type EngineFunc[
	Q documents.Query,
	R any,
] = func(context.Context, Q) (R, error)

const REQUEST_BODY_LIMIT = 1024

func Init() (engine.Engine, error) {
	logs.FmtDefault()
	logs.SetLevel(logs.INFO)

	return kenja2.ConnectAtlasAws(logs.Logger())
}

func Handler[
	Q documents.Query,
	A any,
	R ~*A,
](f EngineFunc[Q, R]) HadlerFunc[Q, R] {
	return func(ctx context.Context, q Q) (R, error) {
		if size := q.Size(); size == 0 || size > REQUEST_BODY_LIMIT {
			err := fmt.Errorf("invalid request size: %d", size)
			logs.Warn(err)
			return nil, err
		}

		res, err := f(ctx, q)
		if err != nil {
			logs.Error(err)
			return nil, err
		}

		return res, nil
	}
}
