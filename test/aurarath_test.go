package test

import (
	"testing"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/aursir4go/aurarath/implementation/zaurarath"
	"github.com/joernweissenborn/stream2go"
	"time"
	"github.com/joernweissenborn/future2go"
)


func TestNodeDiscover(t *testing.T) {

	node1 := aurarath.NewNode()
	node1.RegisterImplementation(new(zaurarath.Implementation))
	c1 := make(chan interface{})
	node1.NewPeers().Listen(testlistener(c1))

	node2 := aurarath.NewNode()
	c2 := make(chan interface{})
	node2.RegisterImplementation(new(zaurarath.Implementation))
	node2.NewPeers().Listen(testlistener(c2))

	node1.Run()
	node2.Run()


	select {
	case <-time.After(2 * time.Second):
		t.Error("Beaon didnt stop")
	case data := <-c1:
		if string(data.(*aurarath.Peer).Id) != string(node2.Self.Id) {
			t.Error("wrong id", data)
		}
		if !node1.KnownPeer(data) {
			t.Error("peer is not nown")
		}

	}


	select {
	case <-time.After(2 * time.Second):
		t.Error("Beaon didnt stop")
	case data := <-c2:
		if string(data.(*aurarath.Peer).Id) != string(node1.Self.Id)  {
			t.Error("wrong id", data)
		}
		if !node2.KnownPeer(data) {
			t.Error("peer is not nown")
		}
	}


	node1.Stop()
	node2.Stop()


}


func TestNodeTimeOut(t *testing.T) {

	node1 := aurarath.NewNode()
	node1.RegisterImplementation(new(zaurarath.Implementation))
	c1 := make(chan interface{})
	node1.LeavingPeers().Listen(testlistener(c1))

	node2 := aurarath.NewNode()
	node2.RegisterImplementation(new(zaurarath.Implementation))

	node1.Run()
	node2.Run()
	time.Sleep(2*time.Second)
	node2.Stop()

	select {
	case <-time.After(10 * time.Second):
		t.Error("Beaon didnt stop")
	case data := <-c1:
		if string(data.(*aurarath.Peer).Id) != string(node2.Self.Id) {
			t.Error("wrong id", data)
		}
		if node1.KnownPeer(data) {
			t.Error("peer is known")
		}

	}



	node1.Stop()



}



func TestNodeConnection(t *testing.T) {

	node1 := aurarath.NewNode()
	node1.RegisterImplementation(new(zaurarath.Implementation))
	c1 := make(chan interface{})
	node1.NewPeers().Listen(testlistener(c1))

	node2 := aurarath.NewNode()
	node2.RegisterImplementation(new(zaurarath.Implementation))

	node1.Run()
	node2.Run()
	select {
	case <-time.After(2 * time.Second):
		t.Error("Peer wasnt found")

	case data := <-c1:
		p :=data.(*aurarath.Peer)

		f := node1.OpenConnection(p)
		c2 := make(chan interface{})
		c3 := make(chan interface{})
		f.Then(testcompleter(c2))
		f.Err(testcompletererr(c3))

		select {
		case <-time.After(5 * time.Second):
			t.Error("Got Nothing")
		case data := <-c3:
				t.Error(" err", data)
		case <-c2:
		}



		node1.Stop()
		node2.Stop()



	}

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
