package aurarath

import "github.com/joernweissenborn/future2go"

type Peer struct {
	Id     []byte
	Codecs []uint8

	Addresses []Address

	Connected *future2go.Future
}
