package zaurarath

import (
	"time"
	"net"
	"log"
	"github.com/joernweissenborn/stream2go"
	"bytes"
)

const (
	BROADCASTADDRESS = "224.0.0.251"
	BROADCASTPORT = 5556
)

type Beacon struct {
	payload []byte
	kill chan struct {}
	autotrack bool
	listensock *net.UDPConn
	outsocks []*net.UDPConn

	in stream2go.StreamController
}

func NewBeacon(payload []byte) (b Beacon) {
	b.Setup()
	b.payload = payload
	b.in = stream2go.New()
	return
}

func (b *Beacon) Setup() {
	b.kill = make(chan struct{})
	b.setupBroadcastListener()
	b.outsocks= []*net.UDPConn{}
	Interfaces,_ := net.Interfaces()
	for _,iface := range Interfaces {

		b.setupBeacon(iface)
	}
}

func (b Beacon) Stop() {
	b.kill <- struct{}{}
}


func (b *Beacon) setupBroadcastListener() (err error) {

	b.listensock, err = net.ListenMulticastUDP("udp4",nil, &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 251),
		Port: BROADCASTPORT,
	})
	return
}



func (b *Beacon) setupBeacon(Interface net.Interface) (err error) {
	BROADCAST_IPv4 := net.IPv4bcast
	ip,_ := Interface.Addrs()
	addr, err := net.ResolveIPAddr("ip4", ip[0].String())
	if err != nil {
		log.Fatal(err)
	}

	socket, err := net.DialUDP("udp4", &net.UDPAddr{
		IP:   addr.IP,
		Port: 0,},
		&net.UDPAddr{
		IP:   BROADCAST_IPv4,
		Port: BROADCASTPORT,
	})

	b.outsocks = append(b.outsocks, socket)

	return
}




func (b Beacon) Run(){
		go b.Listen()
		for _, s := range b.outsocks {
			go b.Ping(s)
		}
}


func (b Beacon) Listen(){


	for {
		data := make([]byte, 1024)
		read, remoteAddr, _ := b.listensock.ReadFromUDP(data)

		b.in.Add(Signal{remoteAddr.IP.String(),data[:read]})
	}

}

func (b Beacon) Signals() stream2go.Stream{
	return b.in.Where(b.noEcho)
}

func (b Beacon) Ping(s *net.UDPConn) {

	var pingtime = 1 * time.Second

	t := time.NewTimer(pingtime)
	for {
		select {
		case <-b.kill:
			return

		case <-t.C:
			s.Write(b.payload)
			t.Reset(pingtime)
		}
	}

}

func (b Beacon) noEcho(d interface {}) bool {
	return !bytes.Equal(d.(Signal).Data,b.payload)
}

type Signal struct {
	SenderIp string
	Data []byte
}

