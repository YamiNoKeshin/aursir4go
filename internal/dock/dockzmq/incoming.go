package dockzmq

import (
	"net"
	"github.com/pebbe/zmq4"
	"strconv"
	"github.com/joernweissenborn/aursir4go/internal/dock"
	"github.com/joernweissenborn/aursir4go/internal"
)

type IncomingZmq struct {
	port int64
	skt *zmq4.Socket
	proc *internal.Incomingprocessor
}

func (izmq IncomingZmq) Activate(proc *internal.Incomingprocessor) (outgoing dock.Outgoing ,err error){
	izmq.port = getRandomPort()
	izmq.skt, err = zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		return
	}

	err = izmq.skt.Bind("tcp://*:" + strconv.FormatInt(izmq.port, 10))
	if err != nil {
		return
	}
	outgoing = izmq.GetOutgoing()
	go izmq.listener()
	return
}
func (IncomingZmq) Deactivate()
func (izmq IncomingZmq) GetOutgoing() OutgoingZmq {
	return OutgoingZmq{nil, izmq.port}
}



func getRandomPort() int64 {
	l, err := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	if err != nil {
		panic("Could not find a free port")
	}
	defer l.Close()
	return int64(l.Addr().(*net.TCPAddr).Port)
}


func (izmq IncomingZmq) listener() {

	defer izmq.skt.Close()

	for {

		msg, err := izmq.skt.RecvMessage(0)

		if err == nil {
			msgtype, _ := strconv.ParseInt(msg[1], 10, 64)
			codec := msg[2]
			encmsg := []byte(msg[3])
			go izmq.proc.ProcessMsg(msgtype,codec,encmsg)

		}

	}


}
