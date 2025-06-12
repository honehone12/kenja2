package mongodb

import (
	"kenja2/endec"
	"kenja2/engine"
)

var _ engine.Engine = &Atlas[endec.MsgPack, endec.Json]{
	encoder: endec.MsgPack{},
	decoder: endec.Json{},
}
