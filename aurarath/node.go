package aurarath

import (
	"github.com/joernweissenborn/stream2go"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"sync"
	"github.com/joernweissenborn/future2go"
	"errors"
	"bytes"
	"time"
)

type Node struct {
	*sync.RWMutex
	Self          *Peer
	newPeers     stream2go.StreamController
	leavingPeers stream2go.StreamController
	in           stream2go.StreamController
	peers        map[string]*Peer
	implementations []Implementation
}

func NewNode() (n *Node) {
	n = new(Node)
	n.Self = NewPeer(generateUuid())
	n.RWMutex = new(sync.RWMutex)
	n.newPeers = stream2go.New()
	n.leavingPeers = stream2go.New()
	n.in = stream2go.New()
	n.implementations = []Implementation{}
	n.peers = map[string]*Peer{}
	return
}

func (n *Node) Run() {
	n.Lock()
	defer  n.Unlock()
	for _, i := range n.implementations {
		i.Init(n.Self.Id())
		i.Run()
		k, u := i.NewPeers().Split(n.KnownPeer)
		u.Listen(n.newPeer)
		k.WhereNot(n.knownPeerAddress).Listen(n.newPeerAddress)
		i.LeavingPeers().Listen(n.removePeerAddress)
		n.in.Join(i.In())
		for _,a := range i.GetAdresses() {
			n.Self.AddAddress(a)
		}
	}
	n.in.Where(isProtocol(PROTOCOL_CONSTANT)).Where(isHello).Listen(n.onHello)

}

func (n *Node) Stop() {
	for _, i := range n.implementations {
		i.Stop()
	}

}

func (n *Node) newPeer(d interface{}){
	p := d.(*Peer)
	n.Lock()
	defer  n.Unlock()
	n.peers[p.IdString()] = p
	n.newPeers.Add(p)


}
var hellos = 0
var hellosc = 0
var hellosn = 0
func (n *Node) onHello(d interface{}){
	hellos++
	log.Println("HELLONR", hellos)
	m := toHello(d)

	if n.KnownPeer(m.Sender) {

		if !n.knownPeerAddress(m.Sender) {
			n.newPeerAddress(m.Sender)
		}

	} else {
		n.newPeer(m.Sender)
	}

	n.Lock()
	defer n.Unlock()
	peer := n.peers[m.Sender.IdString()]
	peer.SetCodecs(m.Codecs)

	if !peer.Connected().IsComplete() {
		hellosn++
		hm := n.encodeMsg(HelloMessage{})
		for _, i:= range n.implementations {
			if is, add := i.Responsible(*peer); is {
				out, err := i.Connect(*n.Self,add)
				if err != nil {
					log.Println("ERROR OPEN CHAN", out)
					defer n.onHello(d)
					return
				}
				out.Add(hm)
				s := stream2go.New()
				out.Join(s.Transform(n.encodeMsg))
				peer.Connected().Complete(s)
			}
			break
		}
	}
	log.Println("HELLONC", hellosc)
	log.Println("HELLONN", hellosn)

	if !peer.Connected().IsComplete() {
		peer.Connected().CompleteError(errors.New("NO_IMP"))
		peer.ResetConnected()
	}

}

func (n *Node) removePeer(p *Peer){
	delete(n.peers,p.IdString())
	n.leavingPeers.Add(p)
}

func (n *Node) removePeerAddress(d interface{}){
	p := d.(LeavingPeerAddress)
	n.Lock()
	defer  n.Unlock()

	peer, f := n.peers[string(p.id)]
	if !f {return}
	peer.RemoveAddress(p.address)
	if len(peer.Addresses())==0 {
		n.removePeer(peer)
	}
}

func (n Node) KnownPeer(d interface{}) (is bool) {

	var id string
	if p, f := d.(Peer); f {
		id = p.IdString()
	} else if p, f := d.(*Peer); f {
		id = p.IdString()
	}
	n.RLock()
	defer n.RUnlock()
	_, is = n.peers[id]
	return
}
func (n Node) newPeerAddress(d interface{}) {
	p := d.(NewPeerAddress)
	n.Lock()
	defer n.Unlock()
	kp := n.peers[string(p.Id)]

	kp.AddAddress(p.Adress)
	return
}

func (n Node) knownPeerAddress(d interface{}) (is bool) {
	p := d.(Peer)
	n.RLock()
	defer n.RUnlock()
	func(d interface{}) interface{} {
		snd++
		log.Println("Sndpeer",snd)

		d.(stream2go.StreamController).Add(m)
		return nil
	}
	for _, a := range n.GetPeer(p).Addresses {
		if a.Implementation == p.Addresses[0].Implementation {
			return true
		}
	}
	return
}
func (n Node) GetPeer(d interface{}) ( *Peer) {
	p := d.(Peer)
	return n.peers[string(p.Id)]
}

func (n *Node) OpenConnection(target *Peer) *future2go.Future {
	if target.Connected != nil {
		return target.Connected
	}
	n.Lock()
	defer n.Unlock()
	target.Connected = future2go.New()
	m := n.encodeMsg(HelloMessage{})
	for _, i:= range n.implementations {
		if is, add := i.Responsible(*target); is {
			out, err := i.Connect(*n.Self,add)
			out.Add(m)
			if err != nil {
				log.Println("ERROR OPEN1 CHAN", out)
				defer n.OpenConnection(target)
				return nil
			}
			time.Sleep(10*time.Millisecond)
			out.Close()
			return target.Connected
		}
	}
	f := target.Connected
	target.Connected = nil
	f.CompleteError(errors.New("NO_IMP"))
	return f
}
func (n *Node) NewPeers() stream2go.Stream {
	return n.newPeers.Stream

}
func (n *Node) encodeMsg(d interface {}) interface {} {

	m := d.(ProtocolMessage)
	var msg Message
	msg.Sender = *n.Self
	msg.Protocol = m.ProtocolId()
	msg.Type= m.Type()
	msg.Payloads = []Payload{}
	for _, d := range m.Data() {
		if b, ok := d.([]byte); ok {
			msg.Payloads = append(msg.Payloads,Payload{BIN, bytes.NewBuffer(b)})
		} else {
			msg.Payloads = append(msg.Payloads,Payload{JSON, encode(d)})
		}
	}
	return msg
}
func (n Node) LeavingPeers() stream2go.Stream {
	return n.leavingPeers.Stream
}

func (n *Node) RegisterImplementation(i Implementation) {

	n.implementations = append(n.implementations,i)
}

func (n *Node) RegisterProtocol(p Protocol) (in stream2go.Stream){
	return n.in.Where(isProtocol(p.ProtocolId()))
}

func isProtocol(id uint8) stream2go.Filter {
	return func(d interface {}) bool {
		m, ok := ToMessage(d)
		if !ok {
			return ok
		}
		return m.Protocol == id
	}
}


func generateUuid() (id []byte) {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return []byte{}
	}
	id = Uuid[0:16]
	if bytes.Equal(id[:1],[]byte{0}) {
		id = generateUuid()
	}
	return
}

func testlistener() stream2go.Suscriber {
	return func(d interface{}) {
		log.Println("BLUB",d)
	}
}
