package aurmellon

import (
	"testing"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/aursir4go/aurarath/implementation/zaurarath"
	"github.com/joernweissenborn/stream2go"
	"time"
	"github.com/joernweissenborn/future2go"
	"log"
)


func TestNodeSupernodeGreeting(t *testing.T) {

	servernode := aurarath.NewNode()
	servernode.RegisterImplementation(new(zaurarath.Implementation))
	server := NewInterface(servernode)
	c1 := make(chan interface{})
	servernode.NewPeers().Listen(testlistener(c1))
	c2 := make(chan interface{})
	server.in.Where(IsHelloIamNodeMessage).Listen(testlistener(c2))
	clientnode := aurarath.NewNode()
	clientnode.RegisterImplementation(new(zaurarath.Implementation))
	client := NewInterface(clientnode)

	servernode.Run()
	clientnode.Run()


	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Peer not found")
	case data := <-c1:
		log.Println("PEER")
		peer := data.(*aurarath.Peer)
		server.send(peer,HelloIamSuperNodeMessage{})
	}


	select {
	case <-time.After(5 * time.Second):
		t.Error("no response")
	case data := <-c2:
		m := data.(aurarath.Message)
		if m.Protocol != PROTOCOL_ID || m.Type != TYPE_HELLO_I_AM_NODE {
			t.Error("wrong msg", data)
		}
		if _, f := client.supernodes[servernode.Self.IdString()];!f {
			t.Error("supernode not registered",string(servernode.Self.Id),client.supernodes)

		}
	}


	servernode.Stop()
	clientnode.Stop()


}

func testlistener(c chan interface{}) stream2go.Suscriber {
	return func(d interface{}) {
		c <- d
	}
}


func testcompleter(c chan interface{}) future2go.CompletionFunc {
	return func(d interface{}) interface {}{
		c <- d
		return nil
	}
}
func testcompletererr(c chan interface{}) future2go.ErrFunc {
	return func(d error) (interface {},error){
		c <- d
		return nil,nil
	}
}
