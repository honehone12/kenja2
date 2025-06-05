package mongodb

import (
	"kenja2"
	"kenja2/marshalers"
)

var _ kenja2.Engine[
	marshalers.MsgPack,
	marshalers.Json,
] = &Atlas[marshalers.MsgPack, marshalers.Json]{
	encoder: marshalers.MsgPack{},
	decoder: marshalers.Json{},
}
