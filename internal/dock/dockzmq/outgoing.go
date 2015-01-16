package dockzmq

import (
	"github.com/pebbe/zmq4"
	"strconv"
)


type OutgoingZmq struct {
	skt *zmq4.Socket
	port int64
	kill chan struct{}
}

func (ozmq *OutgoingZmq) Activate(id string,) (err error){
	ozmq.skt, err = zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		return
	}
	ozmq.skt.SetIdentity(id)
	err = ozmq.skt.Connect("tcp://localhost:5555")
	ozmq.kill = make(chan struct {})
	go pingUdp(id,ozmq.kill)

	return
}
func (ozmq OutgoingZmq)Send(msgtype int64, codec string,msg []byte) (err error){
	ozmq.skt.SendMessage(
		[]string{
		strconv.FormatInt(msgtype, 10),
		codec,
		string(msg),
		strconv.FormatInt(ozmq.port, 10),
		"127.0.0.1",
	}, 0)

	return
}
func (ozmq OutgoingZmq) Close() (err error){
	ozmq.kill <- struct{}{}
	close(ozmq.kill)
	err = ozmq.skt.Close()
	return
}
