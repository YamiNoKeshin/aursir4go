package aurarath

import (
	"github.com/joernweissenborn/stream2go"
)

type Implementation interface {

	Init([]byte) error

	NewPeers() stream2go.Stream
	LeavingPeers() stream2go.Stream
	In() stream2go.Stream

	Responsible(Peer) (bool, Address)

	Connect(home Peer, target Address) (s stream2go.StreamController, err error)

	GetAdresses() (adresses []Address)

	Run()
	Stop()
}
