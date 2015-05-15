package aurarath

type Protocol interface {
	ProtocolId() int8
	Received(m IncomingMessage) interface {}
}
