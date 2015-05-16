package zaurarath

import (
	"github.com/pebbe/zmq4"
	"fmt"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/stream2go"
	"net"
	"encoding/binary"
)


type Outgoing struct {
	skt *zmq4.Socket
	ipportbytes []byte
}

func NewOutgoing(home aurarath.Peer, target aurarath.Address) (out stream2go.StreamController, err error){
	var o Outgoing
	o.skt, err = zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		return
	}

	o.skt.SetIdentity(string(home.Id))
	targetdetails := target.Details.(Details)
	Ip := net.IPv4(uint8(targetdetails.Ip[0]),uint8(targetdetails.Ip[1]),uint8(targetdetails.Ip[2]),uint8(targetdetails.Ip[3]))
	err = o.skt.Connect(fmt.Sprintf("tcp://%s:%d",Ip.String(),targetdetails.Port))

	homedetails, f := FindBestAddress(home,target)
	if !f {
		return
	}
	bp := make([]byte, 2)
	binary.LittleEndian.PutUint16(bp, uint16(homedetails.Details.(Details).Port))
	o.ipportbytes = homedetails.Details.(Details).Ip
	for _, b := range bp {
		o.ipportbytes = append(o.ipportbytes, b)
	}
	out = stream2go.New()
	out.Stream.Listen(o.send)
	out.Stream.Closed.Then(o.Close)
	return
}

func (o Outgoing)send(d interface {}) (){
	msg := ToMessage(d).raw
	msg[2] = o.ipportbytes

	o.skt.SendMessage(msg, 0)

	return
}
func (o Outgoing) Close(interface {}) (interface {}){
	return o.skt.Close()

}
