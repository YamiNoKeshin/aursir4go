package aurmellon

import "github.com/joernweissenborn/aursir4go/aurarath"

const (
	PROTOCOL_ID uint8 = 1
	TYPE_HELLO_I_AM_SUPERNODE        uint8 = 0
	TYPE_HELLO_I_AM_NODE        uint8 = 1
)

type AurMellon struct {

}

func (AurMellon) ProtocolId() uint8 {
	return PROTOCOL_ID
}


type HelloIamSuperNodeMessage struct {
}

func (HelloIamSuperNodeMessage) ProtocolId() uint8 {return PROTOCOL_ID}
func (HelloIamSuperNodeMessage) Type() uint8 {return TYPE_HELLO_I_AM_SUPERNODE}
func (HelloIamSuperNodeMessage) Data() []interface{}{
	return []interface {}{}
}

func IsHelloIamSuperNodeMessage(d interface {}) bool {
	m, ok := aurarath.ToMessage(d)
	if !ok {
		return ok
	}
	return m.Type == TYPE_HELLO_I_AM_SUPERNODE

}

type HelloIamNodeMessage struct {
}

func (HelloIamNodeMessage) ProtocolId() uint8 {return PROTOCOL_ID}
func (HelloIamNodeMessage) Type() uint8 {return TYPE_HELLO_I_AM_NODE}
func (HelloIamNodeMessage) Data() []interface{}{
	return []interface {}{}
}

func IsHelloIamNodeMessage(d interface {}) bool {
	m, ok := aurarath.ToMessage(d)
	if !ok {
		return ok
	}
	return m.Type == TYPE_HELLO_I_AM_NODE

}


