package zaurarath

import (
	"github.com/pebbe/zmq4"
	"fmt"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/stream2go"
)

func outlistener(d interface {}) {

}


type Outgoing struct {
	skt *zmq4.Socket
	ipportbytes []byte
}

func NewOutgoing(home aurarath.Address, target aurarath.Address) (out stream2go.StreamController, err error){
	var o Outgoing
	o.skt, err = zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		return
	}


	//log.Println("ASIp",ip)
	o.skt.SetIdentity(home.Id)
	targetdetails := target.Details.(Details)
	err = o.skt.Connect(fmt.Sprintf("tcp://%s:5555",targetdetails.Ip.String(),targetdetails.Port))

	homedetails := home.Details.(Details)
	o.ipportbytes = append(homedetails.Ip, []byte{homedetails.Port})
	out = stream2go.New()
	out.Stream.Listen(o.send)
	out.Stream.Closed.Then(o.Close)
	return
}

func (o Outgoing)send(d interface {}) (){
	msg := ToMessage(d).raw
	msg[3] = o.ipportbytes
	o.skt.SendMessage(msg, 0)

	return
}
func (o Outgoing) Close(interface {}) (interface {}){
	return o.skt.Close()

}
