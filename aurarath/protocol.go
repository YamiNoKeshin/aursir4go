package aurarath

type Protocol interface {
	ProtocolId() uint8
}

type ProtocolMessage interface {
	ProtocolId() uint8
	Type() uint8
	Data() []interface{}
}
