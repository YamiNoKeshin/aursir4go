package zaurarath

import (
	"github.com/joernweissenborn/stream2go"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"log"
)


type Implementation struct {
	np stream2go.Stream
	in stream2go.Stream

	tracker *PeerTracker
}

func New(uuid []byte) (i Implementation){

	incoming, err := NewIncoming()
	if err != nil {
		log.Fatal(err)
	}
	i.in = incoming.in.Where(MessageOk).Transform(ToIncomingMessage)
	i.tracker = NewPeerTracker()
	go i.tracker.Track()
	var h, l uint8 = uint8(incoming.port>>8), uint8(incoming.port&0xff)

	beaconpayload := []byte{PROTOCOLL_SIGNATURE}

	for _ , b := range uuid {

	}
	beaconpayload := append(uuid,byte(h))
	beaconpayload = append(beaconpayload,byte(l))
	log.Println("Payload",beaconpayload)
	beacon := NewBeacon(beaconpayload)
	k, u := beacon.Signals().Split(i.tracker.isKnown)
	k.Where(i.tracker.isTracked).Listen(i.tracker.Heartbeat)
	u.Listen(i.tracker.add)
	i.np = u
	beacon.Run()

	return
}


func (i Implementation) NewPeers() (s stream2go.Stream) {
	return i.np
}

func (i Implementation) LeavingPeers() (s stream2go.Stream) {
	return i.tracker.deadPeers.Stream
}

func (i Implementation) In() (s stream2go.Stream) {return i.in}

func (i Implementation) RegisterProtocol(p aurarath.Protocol) (s stream2go.Stream) {return}
func (i Implementation) Connect(address aurarath.Address){

}


func filterBeacon(d interface {}) bool {
	return len(d.(Signal)) == 19
}
