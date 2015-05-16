package zaurarath

import (
	"bytes"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"net"
	"encoding/binary"
)

// Message
//
//   0: UUID
//   1: SIGNATURE
//   2: VERSION
//   3: IP+PORT
//   4: PROTOCOL
//   5: TYPE
//   6...: PAYLOAD: CODEC+BYTES
//   CODEC: Byte

const (
	PROTOCOLL_SIGNATURE uint8 = 0xA5
	PROTOCOLL_MAJOR uint8 = 0
	PROTOCOLL_MINOR uint8 = 3
)


type Message struct {
	raw [][]byte
}

func MessageFromRaw(d interface {}) (m interface {}){
	data := [][]byte{}
	for _, di := range d.([]string){
		data = append(data, []byte(di))
	}
	return Message{data}
}

func MessageOk(d interface {}) bool {
	m := ToMessage(d)
	if len(m.raw) < 4 {
		return false
	}

	if !bytes.Equal(m.raw[1], []byte{byte(PROTOCOLL_SIGNATURE)}){
		return false
	}

	if len(m.raw[2]) != 2 {
		return false
	}

	if len(m.raw[3]) != 6 {
		return false
	}
	if len(m.raw[4]) != 1 {
		return false
	}

	if len(m.raw[5]) != 1 {
		return false
	}

	return true
}

func ToIncomingMessage(d interface {}) interface {}{
	data:= d.(Message).raw

	var msg aurarath.Message
	var peer aurarath.Peer
	peer.Id = data[0]
	var address aurarath.Address
	address.Implementation = IMPLEMENTATION_STRING
	var details Details
	rawdetails := data[3]
	details.Ip = net.IPv4(uint8(rawdetails[0]),uint8(rawdetails[1]),uint8(rawdetails[2]),uint8(rawdetails[3]))
	details.Ip = details.Ip[len(details.Ip)-5:]
	details.Port = binary.LittleEndian.Uint16(rawdetails[:4])
	address.Details= details
	msg.Sender = peer

	msg.Protocol = uint8(data[4][0])
	msg.Type = uint8(data[5][0])
	msg.Payloads = []aurarath.Payload{}
	for i := 6; i< len(data);i++ {
		if len(data[i]) <= 1 {
			continue
		}
		msg.Payloads = append(msg.Payloads, aurarath.Payload{uint8(data[i][len(data[i])-1]), bytes.NewBuffer(data[i][:len(data[i])-1])})
	}
	return msg
}
func OutgoingToMessage(d interface {}) interface {}{
	m := d.(aurarath.Message)

	var msg Message
	msg.raw = [][]byte{
		[]byte{byte(PROTOCOLL_SIGNATURE)},
		[]byte{byte(PROTOCOLL_MAJOR),byte(PROTOCOLL_MINOR)},
		[]byte{},
		[]byte{byte(m.Protocol)},
		[]byte{byte(m.Type)},
	}

	for _,pl := range m.Payloads {
		load := append(pl.Bytes.Bytes(),byte(pl.Codec))
		msg.raw = append(msg.raw,load)
	}

	return msg
}

func ToMessage(d interface {}) (m Message){
	m, _ = d.(Message)
	return
}
