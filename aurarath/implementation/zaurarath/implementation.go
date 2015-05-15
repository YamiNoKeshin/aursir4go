package zaurarath

import (
	"github.com/joernweissenborn/stream2go"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"log"
)


type Implementation struct {
	np stream2go.StreamController
	lp stream2go.StreamController
	in stream2go.Stream

	tracker *PeerTracker
}

func New(uuid []byte) (i Implementation){

	incoming, err := NewIncoming()
	if err != nil {
		log.Fatal(err)
	}
	i.in = incoming.in.Where(MessageOk).Transform(ToIncomingMessage)
	i.tracker := NewPeerTracker()
	go i.tracker.Track()

	beaconpayload := append(uuid,[]byte(incoming.port))
	beacon := NewBeacon(beaconpayload)
	beacon.Signals().Where(i.tracker.isTracked).Listen(i.tracker.Heartbeat)
	beacon.Run()

	return
}


func (i Implementation) NewPeers() (s stream2go.Stream) {return}
func (i Implementation) LeavingPeers() (s stream2go.Stream) {return}

func (i Implementation) In() (s stream2go.Stream) {return i.in}

func (i Implementation) RegisterProtocol(p aurarath.Protocol) (s stream2go.Stream) {return}
func (i Implementation) Connect(address aurarath.Address){

}


