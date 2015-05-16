package zaurarath

import (
	"github.com/joernweissenborn/aursir4go/aurarath"
	"bytes"
)

const (
	IMPLEMENTATION_STRING = "ZAURARATH_0_3"
)

type Details struct {

	Ip []byte
	Port uint16

}

func FindBestAddress(peer aurarath.Peer, target aurarath.Address) (match aurarath.Address, f bool){
	if len(peer.Addresses) == 1 {
		return peer.Addresses[1], true
	}
	td := target.Details.(Details)
	for _, addr := range peer.Addresses {
		ad := addr.Details.(Details)
		if bytes.Equal(td.Ip[:2],ad.Ip[:2]) {
			return addr,true
		}
	}
	return
}
