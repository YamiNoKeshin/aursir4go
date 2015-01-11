package dockzmq

import (
	"github.com/pebbe/zmq4"
	"strconv"
)


type OutgoingZmq struct {
	skt *zmq4.Socket
	port int64
}

func (ozmq *OutgoingZmq) Activate(id string,) (err error){
	ozmq.skt, err = zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		return
	}
	ozmq.skt.SetIdentity(id)
	err = ozmq.skt.Connect("tcp://localhost:5555")

	return
}
func (ozmq OutgoingZmq)Send(msgtype int64, codec string,msg []byte) (err error){
	ozmq.skt.SendMessage(
		[]string{
		strconv.FormatInt(msgtype, 10),
		codec,
		string(msg),
		strconv.FormatInt(ozmq.port, 10)}, 0)

	return
}
