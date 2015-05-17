package zaurarath

import (
	"bytes"
	"encoding/binary"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/stream2go"
	"log"
	"net"
)

type Implementation struct {
	np      stream2go.Stream
	in      stream2go.Stream
	b       Beacon
	tracker *PeerTracker
	id []byte
	port    uint16
}

func New(uuid []byte) (i Implementation) {

	i.Init(uuid)
	i.Run()
	return
}

func (i *Implementation) Responsible(p aurarath.Peer) (bool, aurarath.Address){
	return i.tracker.isKnownAddr(p)
}
func (i *Implementation) Init(UUID[]byte) error{
	i.id = UUID
	return nil
}
func (i *Implementation) Run(){
	incoming, err := NewIncoming()
	if err != nil {
		log.Fatal(err)
	}
	i.in = incoming.in.Transform(MessageFromRaw).Where(MessageOk).Transform(ToIncomingMessage)
	i.tracker = NewPeerTracker()
	go i.tracker.Track()
	beaconpayload := []byte{PROTOCOLL_SIGNATURE}

	for _, b := range i.id {
		beaconpayload = append(beaconpayload, byte(b))
	}
	i.port = incoming.port
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i.port)
	beaconpayload = append(beaconpayload, b[0])
	beaconpayload = append(beaconpayload, b[1])
	i.b = NewBeacon(beaconpayload, 5558)
	k, u := i.b.Signals().Where(filterBeacon(i.id)).Transform(parseBeacon).Split(i.tracker.isKnown)
	k.Listen(i.tracker.Heartbeat)
	u.Listen(i.tracker.add)
	i.in.Transform(func(d interface {})interface {}{
		return d.(aurarath.Message).Sender
	}).WhereNot(i.tracker.isKnown).Listen(i.tracker.add)
	i.np = u
	i.b.Run()
}
func (i Implementation) NewPeers() (s stream2go.Stream) {
	return i.np
}

func (i Implementation) LeavingPeers() (s stream2go.Stream) {
	return i.tracker.deadPeers.Stream
}

func (i Implementation) In() (s stream2go.Stream) { return i.in }

func (i Implementation) Connect(home aurarath.Peer, target aurarath.Address) (s stream2go.StreamController, err error) {
	outstream, err := NewOutgoing(home, target)
	if err != nil {
		return
	}
	s = stream2go.New()
	outstream.Join(s.Transform(OutgoingToMessage))
	return
}

func (i Implementation) GetAdresses() (adresses []aurarath.Address) {

	adresses = []aurarath.Address{}
	Interfaces, _ := net.Interfaces()
	for _, iface := range Interfaces {
		var a aurarath.Address
		a.Implementation = IMPLEMENTATION_STRING
		var d Details
		add, _ := iface.Addrs()

		Ip, _ := net.ResolveIPAddr(add[0].Network(), add[0].String())
		d.Ip = Ip.IP[len(Ip.IP)-4:]
		d.Port = i.port
		a.Details = d
		adresses = append(adresses, a)
	}
	return
}

func (i Implementation) Stop() {
	i.np.Close()
	i.in.Close()
	i.b.Stop()
}

func filterBeacon(uuid []byte) func(d interface{}) bool {
	return func(d interface{}) bool {
		if len(d.(Signal).Data) != 19 {
			return false
		}
		return !bytes.Equal(uuid, d.(Signal).Data[1:17])
	}
}

func parseBeacon(d interface{}) interface{} {

	data := d.(Signal).Data
	var p aurarath.Peer
	p.Id = data[1:17]
	var a aurarath.Address
	a.Implementation = IMPLEMENTATION_STRING
	var det Details
	det.Ip = d.(Signal).SenderIp
	det.Port = binary.LittleEndian.Uint16(data[17:])
	a.Details = det
	p.Addresses = []aurarath.Address{a}
	return p
}
