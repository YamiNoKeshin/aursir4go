package dockzmq

import (
	"net"
	"github.com/pebbe/zmq4"
	"strconv"
	"github.com/joernweissenborn/aursir4go/internal/dock"
	"github.com/joernweissenborn/aursir4go/internal"
	"time"
	"fmt"
)

type IncomingZmq struct {
	port int64
	skt *zmq4.Socket
	proc *internal.Incomingprocessor
}

func (izmq IncomingZmq) Activate(proc *internal.Incomingprocessor) (outgoing dock.Outgoing ,err error){
	izmq.port = getRandomPort()
	izmq.proc = proc
	izmq.skt, err = zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		return
	}

	err = izmq.skt.Bind("tcp://*:" + strconv.FormatInt(izmq.port, 10))
	if err != nil {
		return
	}
	o:= izmq.GetOutgoing()

	outgoing = &o
	go izmq.listener()

		return
}
func (IncomingZmq) Deactivate() {
	//TODO Deactivate interface
}
func (izmq IncomingZmq) GetOutgoing() OutgoingZmq {
	return OutgoingZmq{nil, izmq.port,nil,""}
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
		                   //log.Println(msg)
		if err == nil {
			msgtype, _ := strconv.ParseInt(msg[1], 10, 64)
			codec := msg[2]
			encmsg := []byte(msg[3])
			go izmq.proc.ProcessMsg(msgtype,codec,encmsg)

		}

	}


}

func pingUdp(UUID string, kill chan struct {}) {

	var pingtime = 1 * time.Second

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:0","127.0.0.1"))
	if err != nil {
		panic(err)
	}
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5557")

	if err != nil {
		panic(err)
	}
	con, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		panic(err)
	}
	t := time.NewTimer(pingtime)
	fmt.Println(fmt.Sprintf("Beginning UDP Broadcast with %s",UUID))
	for {
		select {
		case <-kill:
			return

		case <-t.C:
			con.Write([]byte(fmt.Sprintf("%s", UUID)))
			t.Reset(pingtime)
		}
	}

}
