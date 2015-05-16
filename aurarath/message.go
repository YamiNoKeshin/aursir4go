package aurarath

import "bytes"

func NewMessage(Sender Peer, Protocol uint8 ,Type uint8, Payloads []Payload) Message {
	return Message{Sender , Protocol ,Type , Payloads}
}

type Message struct {
	Sender Peer
	Protocol uint8
	Type uint8
	Payloads []Payload
}

func (m *Message) Pop(target interface {})  {
//	d := m.Payloads[0]
//	decode(d.Bytes,target,d.Codec)
	m.Payloads = m.Payloads[1:]
}

type Payload struct {
	Codec uint8
	Bytes *bytes.Buffer
}

