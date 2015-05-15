package aurarath

import "github.com/joernweissenborn/stream2go"

type Implementation interface {
	IsProtocol(msg interface {}) bool
	NewPeers() stream2go.Stream
	LeavingPeers() stream2go.Stream
	In() stream2go.Stream
	RegisterProtocol(p Protocol) stream2go.Stream
	Send(interface {})
}
