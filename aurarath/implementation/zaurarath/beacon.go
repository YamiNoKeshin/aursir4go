package zaurarath

import (
	"bytes"
	"github.com/joernweissenborn/future2go"
	"github.com/joernweissenborn/stream2go"
	"log"
	"net"
	"time"
)

const (
	BROADCASTADDRESS = "224.0.0.251"
)

type Beacon struct {
	payload    []byte
	kill       *future2go.Future
	autotrack  bool
	listensock *net.UDPConn
	outsocks   []*net.UDPConn
	port       int
	in         stream2go.StreamController
}

func NewBeacon(payload []byte, port uint16) (b Beacon) {
	b.payload = payload
	b.in = stream2go.New()
	b.port = int(port)
	b.Setup()
	return
}

func (b *Beacon) Setup() {
	b.kill = future2go.New()
	b.setupBroadcastlistener()
	b.outsocks = []*net.UDPConn{}
	Interfaces, _ := net.Interfaces()
	for _, iface := range Interfaces {

		b.setupBeacon(iface)
	}
}

func (b Beacon) Stop() {
	b.kill.Complete(nil)

}

func (b *Beacon) setupBroadcastlistener() (err error) {

	b.listensock, err = net.ListenMulticastUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 251),
		Port: b.port,
	})
	return
}

func (b *Beacon) setupBeacon(Interface net.Interface) (err error) {
	BROADCAST_IPv4 := net.IPv4bcast
	ip, _ := Interface.Addrs()
	addr, err := net.ResolveIPAddr("ip4", ip[0].String())
	if err != nil {
		log.Fatal(err)
	}

	socket, err := net.DialUDP("udp4", &net.UDPAddr{
		IP:   addr.IP,
		Port: 0},
		&net.UDPAddr{
			IP:   BROADCAST_IPv4,
			Port: b.port,
		})

	b.outsocks = append(b.outsocks, socket)

	return
}

func (b Beacon) Run() {
	go b.listen()
	for _, s := range b.outsocks {
		go b.Ping(s)
	}
}

func (b Beacon) listen() {

	c := make(chan struct{})
	kill := b.kill.AsChan()
	go b.getSignal(c)
	for {
		select {
		case <-kill:
			return

		case <-c:
			go b.getSignal(c)
		}
	}

}

func (b Beacon) getSignal(c chan struct{}) {
	data := make([]byte, 1024)
	read, remoteAddr, _ := b.listensock.ReadFromUDP(data)

	b.in.Add(Signal{remoteAddr.IP[len(remoteAddr.IP)-4:], data[:read]})
	c <- struct{}{}
}

func (b Beacon) Signals() stream2go.Stream {
	return b.in.Where(b.noEcho)
}

func (b Beacon) Ping(s *net.UDPConn) {

	var pingtime = 1000 * time.Millisecond
	kill := b.kill.AsChan()

	t := time.NewTimer(pingtime)
	//s.Write(b.payload)
	for {
		select {
		case <-kill:
			return

		case <-t.C:
			s.Write(b.payload)
			t.Reset(pingtime)
		}
	}

}

func (b Beacon) noEcho(d interface{}) bool {
	return !bytes.Equal(d.(Signal).Data, b.payload)
}

type Signal struct {
	SenderIp []byte
	Data     []byte
}
