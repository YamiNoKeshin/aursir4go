package aurarath

import (
	"github.com/joernweissenborn/stream2go"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"sync"
	"github.com/joernweissenborn/future2go"
	"errors"
	"time"
)

type Node struct {
	*sync.RWMutex
	Self           Peer
	newPeers     stream2go.StreamController
	leavingPeers stream2go.StreamController
	in           stream2go.StreamController
	peers        map[string]Peer
	implementations []Implementation
}

func NewNode() (n *Node) {
	n = new(Node)
	n.Self.Id = generateUuid()
	n.Self.Addresses = []Address{}
	n.RWMutex = new(sync.RWMutex)
	n.newPeers = stream2go.New()
	n.leavingPeers = stream2go.New()
	n.in = stream2go.New()
	n.implementations = []Implementation{}
	n.peers = map[string]Peer{}
	return
}

func (n *Node) Run() {
	n.Lock()
	defer  n.Unlock()
	for _, i := range n.implementations {
		i.Init(n.Self.Id)
		i.Run()
		k, u := i.NewPeers().Split(n.KnownPeer)
		n.newPeers.Join(u)
		u.Listen(n.newPeer)
		k.WhereNot(n.knownPeerAddress).Listen(n.newPeerAddress)
		i.LeavingPeers().Listen(n.removePeerAddress)
		n.in.Join(i.In())
		for _,a := range i.GetAdresses() {
			n.Self.Addresses = append(n.Self.Addresses,a)
		}
	}
	n.in.Where(isAurArath).Where(isHello).Listen(n.onHello)

}

func (n *Node) Stop() {
	for _, i := range n.implementations {
		i.Stop()
	}

}

func (n *Node) newPeer(d interface{}){
	p := d.(Peer)
	n.Lock()
	defer  n.Unlock()
	n.peers[string(p.Id)] = p

}

func (n *Node) onHello(d interface{}){
	m := toHello(d)
	peer := m.Sender
	peer.Codecs = m.Codecs
	if n.KnownPeer(peer) {
		if !n.knownPeerAddress(peer) {
			n.newPeerAddress(peer)
		}
		n.Lock()
		peer.Addresses = n.peers[string(peer.Id)].Addresses
		peer.Connected= n.peers[string(peer.Id)].Connected
	} else {
		n.Lock()
	}
	defer n.Unlock()

	if peer.Connected == nil {
		peer.Connected = future2go.New()
	}
	hm := NewHelloMessage(n.Self)

	for _, i:= range n.implementations {
		if is, add := i.Responsible(peer); is {
			out, _ := i.Connect(n.Self,add)
			out.Add(hm)
			if !peer.Connected.IsComplete(){
				peer.Connected.Complete(out)
			} else {
				out.Close()
			}
			break
		}

	}
	if !peer.Connected.IsComplete() {
		peer.Connected.CompleteError(errors.New("NO_IMP"))
		peer.Connected = nil
	}
	n.peers[string(peer.Id)] = peer

}

func (n *Node) removePeer(p Peer){
	delete(n.peers,string(p.Id))
	n.leavingPeers.Add(p)
}

func (n *Node) removePeerAddress(d interface{}){
	p := d.(Peer)
	n.Lock()
	defer  n.Unlock()
	var j int = -1
	peer := n.peers[string(p.Id)]
	for i, a := range peer.Addresses {
		if a.Implementation == p.Addresses[0].Implementation {
			j = i
			break
		}
	}
	if j != -1 {
		if len(peer.Addresses)==1 {
			n.removePeer(peer)
		} else {
			peer.Addresses[j] = peer.Addresses[len(peer.Addresses)-1]
			peer.Addresses= peer.Addresses[:len(peer.Addresses)-2]
		}
		if len(peer.Addresses)==0 {
			n.removePeer(peer)
		}
	}
}

func (n Node) KnownPeer(d interface{}) (is bool) {
	p := d.(Peer)
	n.RLock()
	defer n.RUnlock()
	_, is = n.peers[string(p.Id)]
	return
}
func (n Node) newPeerAddress(d interface{}) {
	p := d.(Peer)
	n.Lock()
	defer n.Unlock()
	kp := n.peers[string(p.Id)]
	kp.Addresses = append(n.peers[string(p.Id)].Addresses,p.Addresses[0])
	n.peers[string(p.Id)] = kp
	return
}
func (n Node) peerConnectionInitialized(p Peer) bool {
	return n.getPeer(p).Connected != nil
}

func (n Node) knownPeerAddress(d interface{}) (is bool) {
	p := d.(Peer)
	n.RLock()
	defer n.RUnlock()
	for _, a := range n.getPeer(p).Addresses {
		if a.Implementation == p.Addresses[0].Implementation {
			return true
		}
	}
	return
}
func (n Node) getPeer(d interface{}) (p Peer) {
	p = d.(Peer)
	return n.peers[string(p.Id)]
}

func (n *Node) OpenConnection(target Peer) *future2go.Future {
	if n.peerConnectionInitialized(target) {
		return n.getPeer(target).Connected
	}
	n.Lock()
	defer n.Unlock()
	kp := n.peers[string(target.Id)]
	kp.Connected = future2go.New()
	n.peers[string(target.Id)] = kp
	m := NewHelloMessage(n.Self)
	for _, i:= range n.implementations {
		if is, add := i.Responsible(target); is {
			out, _ := i.Connect(n.Self,add)
			out.Add(m)
			time.Sleep(10*time.Millisecond)
			out.Close()
			return kp.Connected
		}
	}
	f := kp.Connected
	kp.Connected = nil
	f.CompleteError(errors.New("NO_IMP"))
	return f
}
func (n *Node) NewPeers() stream2go.Stream {
	return n.newPeers.Stream

}

func (n Node) LeavingPeers() stream2go.Stream {
	return n.leavingPeers.Stream
}

func (n *Node) RegisterImplementation(i Implementation) {

	n.implementations = append(n.implementations,i)
}

func generateUuid() []byte {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return []byte{}
	}
	return Uuid[0:16]
}

func testlistener() stream2go.Suscriber {
	return func(d interface{}) {
		log.Println("BLUB",d)
	}
}
