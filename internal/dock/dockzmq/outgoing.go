package dockzmq

import (
	"github.com/pebbe/zmq4"
	"strconv"
	"os"
	"fmt"
	"log"
	"net"
	"strings"
)


type OutgoingZmq struct {
	skt *zmq4.Socket
	port int64
	kill chan struct{}
	myip string
}

func (ozmq *OutgoingZmq) Activate(id string,) (err error){
	ozmq.skt, err = zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		return
	}
	ip := "127.0.0.1"
	ozmq.myip = "172.17.42.1"
	if envip := os.Getenv("AURSIR_RT_PORT"); envip != "" {
		ip = strings.SplitAfterN(strings.Split(envip,":")[1],"",3)[2]
		myip, err := net.ResolveIPAddr("ip4",os.Getenv("HOSTNAME"))
		if err != nil {
			panic(err)
		}
		ozmq.myip = myip.IP.String()
	} else if envip := os.Getenv("AURSIR_RT_IP"); envip != "" {
		ip = envip
		myip, err := net.ResolveIPAddr("ip4",os.Getenv("HOSTNAME"))
		if err != nil {
			panic(err)
		}
		ozmq.myip = myip.IP.String()

	}
	log.Println("ASIp",ip)
	ozmq.skt.SetIdentity(id)
	err = ozmq.skt.Connect(fmt.Sprintf("tcp://%s:5555",ip))
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
		ozmq.myip,
	}, 0)

	return
}
func (ozmq OutgoingZmq) Close() (err error){
	ozmq.kill <- struct{}{}
	close(ozmq.kill)
	err = ozmq.skt.Close()
	return
}
