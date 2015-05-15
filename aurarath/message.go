package aurarath

import "bytes"

type Message struct {
	Protocol int8
	Payloads []Payload
}

func (m *Message) Pop(target interface {})  {
	d := m.Payloads[0]
	decode(d.Bytes,target,d.Codec)
	m.Payloads = m.Payloads[1:]
}

type Payload struct {
	Codec string
	Bytes *bytes.Buffer
}

type IncomingMessage struct {
	Message
	Sender Address
}

type OutgoingMessage struct {
	Message
	Receiver Address
}
