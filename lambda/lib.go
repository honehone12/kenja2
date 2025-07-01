package lambda

import (
	"context"
	"errors"
	"kenja2"
	"kenja2/engine"
	"kenja2/lambda/logs"
)

type HadlerFunc = func(context.Context, []byte) ([]byte, error)
type EngineFunc = func(context.Context, []byte) ([]byte, error)

const REQUEST_BODY_LIMIT = 1024

func Init() (engine.Engine, error) {
	logs.FmtDefault()
	logs.SetLevel(logs.INFO)

	return kenja2.ConnectAtlasAws(
		kenja2.NewMsgPack(),
		kenja2.NewMsgPack(),
		logs.Logger(),
	)
}

func Handler(f EngineFunc) HadlerFunc {
	return func(ctx context.Context, req []byte) ([]byte, error) {
		if len(req) > REQUEST_BODY_LIMIT {
			emsg := "request body size over limit"
			logs.Warn(emsg)
			return nil, errors.New(emsg)
		}

		res, err := f(ctx, req)
		if err != nil {
			logs.Error(err)
			return nil, err
		}

		return res, nil
	}
}
