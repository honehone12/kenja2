package ed

var _ Encoder = Json{}
var _ Decoder = Json{}
var _ Encoder = MsgPack{}
var _ Decoder = MsgPack{}
