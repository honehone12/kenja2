package mongodb

import (
	"kenja2/endec"
	"kenja2/engine"
)

var _ engine.Engine[endec.MsgPack, endec.Json] = &Atlas[
	endec.MsgPack,
	endec.Json,
]{
	encoder: endec.MsgPack{},
	decoder: endec.Json{},
}
