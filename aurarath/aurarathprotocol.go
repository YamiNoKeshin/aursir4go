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

func (HelloMessage) ProtocolId() uint8 {return PROTOCOL_CONSTANT}
func (HelloMessage) Type() uint8 {return TYPE_HELLO}
func (HelloMessage) Data() []interface{}{
	cs := []byte{}
	for _, c := range CODECS {
		cs = append(cs,byte(c))
	}
	return []interface {}{cs}
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
	data := new(bytes.Buffer)
	m.PopBytes(data)
	for _, b:= range data.Bytes(){
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

