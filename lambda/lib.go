package lambda

import (
	"context"
	"fmt"
	"kenja2"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/lambda/logs"
)

type HadlerFunc[Q documents.Query] = func(context.Context, Q) ([]byte, error)
type EngineFunc[Q documents.Query] = func(context.Context, Q) ([]byte, error)

const REQUEST_BODY_LIMIT = 1024

func Init() (engine.Engine, error) {
	logs.FmtDefault()
	logs.SetLevel(logs.INFO)

	return kenja2.ConnectAtlasAws(
		kenja2.NewMsgPack(),
		kenja2.NewJson(),
		logs.Logger(),
	)
}

func Handler[Q documents.Query](f EngineFunc[Q]) HadlerFunc[Q] {
	return func(ctx context.Context, q Q) ([]byte, error) {
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
