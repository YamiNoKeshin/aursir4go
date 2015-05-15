package aurarath

import (
	"github.com/joernweissenborn/stream2go"
	"log"
	uuid "github.com/nu7hatch/gouuid"
)


type Node struct {
	Id string
	newPeers stream2go.StreamController
	leavingPeers stream2go.StreamController
	in stream2go.StreamController
	out stream2go.StreamController
}

func NewNode() (n Node){
	n.Id = generateUuid()
	n.newPeers = stream2go.New()
	n.leavingPeers = stream2go.New()
	n.in = stream2go.New()
	n.out = stream2go.New()
	return
}

func (n Node) NewPeers() stream2go.Stream{
	return n.newPeers.Stream

}


func (n Node) LeavingPeers() stream2go.Stream{
	return n.leavingPeers.Stream
}



func (n Node) RegisterImplementation(i Implementation){
	n.newPeers.Join(i.NewPeers())
	n.leavingPeers.Join(i.NewPeers())
	n.out.Where(i.Responsible).Listen(i.Send)
	n.in.Join(i.In())
}




func generateUuid() string {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return ""
	}
	return Uuid.String()
}
