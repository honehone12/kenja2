package mongodb

import (
	"kenja2"
	"kenja2/endec"
)

var _ kenja2.Engine[endec.MsgPack, endec.Json] = &Atlas[
	endec.MsgPack,
	endec.Json,
]{
	encoder: endec.MsgPack{},
	decoder: endec.Json{},
}
