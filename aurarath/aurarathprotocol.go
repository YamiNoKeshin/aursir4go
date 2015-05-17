package aurarath

import "bytes"

const (
	PROTOCOL_CONSTANT uint8 = 0
	TYPE_HELLO        uint8 = 0
	TYPE_LEAVE        uint8 = 1
)

type AurArathProtocol struct {
}

func (AurArathProtocol) ProtocolId() uint8 {
	return PROTOCOL_CONSTANT
}

type HelloMessage struct {
	Sender Peer
	Codecs []uint8
}

func NewHelloMessage(Home Peer) Message{
	cs := []byte{}
	for _, c := range CODECS {
		cs = append(cs,byte(c))
	}
	p := Payload{BIN,bytes.NewBuffer(cs)}
	return NewMessage(Home,PROTOCOL_CONSTANT,TYPE_HELLO,[]Payload{p})
}

func isAurArath(d interface {}) bool {
	m, ok := ToMessage(d)
	if !ok {
		return ok
	}
	return m.Protocol == PROTOCOL_CONSTANT

}
func isHello(d interface {}) bool {
	m, ok := ToMessage(d)
	if !ok {
		return ok
	}
	return m.Type == TYPE_HELLO

}

func toHello(d interface {}) HelloMessage {
	m, _ := ToMessage(d)
	var hm HelloMessage
	hm.Sender = m.Sender
	hm.Codecs = []uint8{}
	for _, b:= range m.Payloads[0].Bytes.Bytes() {
		hm.Codecs = append(hm.Codecs,uint8(b))
	}
	return hm

}



func isLeave(d interface {}) bool {
	m, ok := ToMessage(d)
	if !ok {
		return ok
	}
	return m.Type == TYPE_LEAVE

}

