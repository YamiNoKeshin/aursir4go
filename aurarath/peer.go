package aurarath

import (
	"github.com/joernweissenborn/future2go"
	"sync"
)

type NewPeerAddress struct {
	Id []byte
	Address Address
}

type LeavingPeerAddress struct {
	Id []byte
	Address Address
}

type Peer struct {
	sync.RWMutex
	id     []byte
	codecs []uint8

	addresses []Address

	connected *future2go.Future

}

func NewPeer(id []byte)(p *Peer){
	p =new(Peer)
	p.RWMutex = new(sync.RWMutex)
	p.addresses = []Address{}
	p.connected = future2go.New()
	return
}

func (p *Peer) IdString() string {
	return string(p.id)
}

func (p *Peer) Id() []byte {
	return p.id
}
func (p *Peer) SetCodecs(c []byte) []byte {
	p.Lock()
	defer p.Unlock()
	p.codecs = c
	return
}

func (p *Peer) Addresses() []Address {
	p.RLock()
	defer p.RUnlock()
	return p.addresses
}
func (p *Peer) AddAddresses(as []Address) {
	for _, a := range as {
		p.AddAddress(a)
	}
	return
}
func (p *Peer) AddAddress(a Address) {
	p.Lock()
	defer p.Unlock()
	p.addresses = append(p.addresses,a)
	return
}

func (p *Peer) RemoveAddress(add Address) {
	p.Lock()
	defer p.Unlock()
	var j int = -1
	for i, a := range p.Addresses() {
		if a.Implementation == add.Implementation {
			j = i
			break
		}
	}
	if j != -1 {
		if len(p.addresses)==1 {
			p.addresses = []Address{}
		} else {
			p.addresses[j] = p.addresses[len(p.addresses)-1]
			p.addresses = p.addresses[:len(p.addresses)-2]
		}

	}
	return
}

func (p *Peer) Connected() *future2go.Future {
	return p.connected
}


func (p *Peer) ResetConnected() *future2go.Future {
	p.Lock()
	defer p.Unlock()
	p.connected = future2go.New()
	return p.connected
}

