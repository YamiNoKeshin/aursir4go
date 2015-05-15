package aurarath

import (
	"github.com/joernweissenborn/stream2go"
	"github.com/joernweissenborn/future2go"
)

type Implementation interface {

	NewPeers() stream2go.Stream
	LeavingPeers() stream2go.Stream

	RegisterProtocol(p Protocol) stream2go.Stream
	Responsible(interface {}) bool

	Connect(Peer) (out stream2go.StreamController, gone future2go.Future)
}
