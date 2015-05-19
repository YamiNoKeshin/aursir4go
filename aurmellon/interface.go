package aurmellon

import (
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/stream2go"
	"sync"
	"github.com/joernweissenborn/future2go"
)

type Interface struct {
	*sync.RWMutex
	node *aurarath.Node

	in stream2go.Stream
	supernodes map[string]*aurarath.Peer
	SupernodeConnected *future2go.Future
	SupernodeDisconnected *future2go.Future

}

func NewInterface(node *aurarath.Node)(i *Interface){
	i = new(Interface)
	i.node = node
	i.SupernodeConnected = future2go.New()
	i.SupernodeDisconnected = future2go.New()
	i.RWMutex = new(sync.RWMutex)

	i.in = node.RegisterProtocol(AurMellon{})
	i.in.Where(IsHelloIamSuperNodeMessage).Listen(i.newSupernode)

	i.supernodes = map[string]*aurarath.Peer{}
	return
}
func (i *Interface) IsSupernodeConnected() bool {
	return len(i.supernodes)!=0
}

func (i *Interface) newSupernode(d interface {}) {
	m, _ := aurarath.ToMessage(d)
	i.Lock()
	defer i.Unlock()

	peer := i.node.GetPeer(m.Sender)
	if len(i.supernodes) == 0 {
		i.SupernodeConnected.Complete(nil)
	}
	i.supernodes[m.Sender.IdString()] = peer
	i.send(peer,HelloIamNodeMessage{})
}

func (i *Interface) send(p *aurarath.Peer, m aurarath.ProtocolMessage) {
	i.node.OpenConnection(p).Then(func(d interface {})interface {}{

		d.(stream2go.StreamController).Add(m)
		return nil
	})
}
