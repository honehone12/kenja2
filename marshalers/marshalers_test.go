package marshalers

var _ Marshaler = Json{}
var _ Marshaler = Bson{}
var _ Marshaler = MsgPack{}
