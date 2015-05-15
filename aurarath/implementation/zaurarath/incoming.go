package zaurarath

import (
	"net"
	"github.com/pebbe/zmq4"
	"strconv"
	"github.com/joernweissenborn/stream2go"
)

type Incoming struct {
	port uint8
	skt *zmq4.Socket
	in stream2go.StreamController
}


func NewIncoming() (i Incoming, err error){
	i.in = stream2go.New()
	err = i.setupSocket()
	return
}
func (i *Incoming) setupSocket() (err error){
	i.port = getRandomPort()
	i.skt, err = zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		return
	}

	err = i.skt.Bind("tcp://*:" + strconv.FormatInt(i.port, 10))
		return
}

func getRandomPort() int64 {
	l, err := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	if err != nil {
		panic("Could not find a free port")
	}
	defer l.Close()
	return int64(l.Addr().(*net.TCPAddr).Port)
}


func (i Incoming) Listen() {

	defer i.skt.Close()

	for {

		msg, err := i.skt.RecvMessage(0)
		                   //log.Println(msg)
		if err == nil {
			i.in.Add(msg)
		}
	}

}
