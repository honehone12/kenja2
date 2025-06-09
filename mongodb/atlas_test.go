package mongodb

import (
	"kenja2"
	"kenja2/ed"
)

var _ kenja2.Engine[ed.MsgPack, ed.Json] = &Atlas[
	ed.MsgPack,
	ed.Json,
]{
	encoder: ed.MsgPack{},
	decoder: ed.Json{},
}
