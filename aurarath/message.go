package aurarath

import (
	"bytes"
	"io"
)

func NewMessage(Sender Peer, Protocol uint8, Type uint8, Payloads []Payload) Message {
	return Message{Sender, Protocol, Type, Payloads}
}

type Message struct {
	Sender   Peer
	Protocol uint8
	Type     uint8
	Payloads []Payload
}

func (m *Message) Pop(target interface{}) bool {
	if len(m.Payloads)==0 {
		return false
	}
	d := m.Payloads[0]
	decode(d.Bytes,&target,d.Codec)
	m.Payloads = m.Payloads[1:]
	return len(m.Payloads)!=0
}

func (m *Message) PopBytes(d *bytes.Buffer) bool {
	if len(m.Payloads)==0 {
		return false
	}
	io.Copy(d, m.Payloads[0].Bytes)

	return len(m.Payloads)!=0
}

type Payload struct {
	Codec uint8
	Bytes *bytes.Buffer
}

func ToMessage(d interface {})(m Message, ok bool){
	m, ok = d.(Message)
	return
}
