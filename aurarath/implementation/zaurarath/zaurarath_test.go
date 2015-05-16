package zaurarath

import (
	"testing"
	"github.com/joernweissenborn/stream2go"
	"log"
)


func TestPeerDiscover(t *testing.T){

	iface1 := New([]byte("1123421253"))
	c1 := make(chan interface {})
	iface1.NewPeers().Listen(testlistener(c1))

	iface2 := New([]byte("2243566326342"))
	c2 := make(chan interface {})
	iface2.NewPeers().Listen(testlistener(c2))
	data := string((<-c1).(Signal).Data)
	log.Println(data)
	log.Println(<-c2)

}

func testlistener(c chan interface {}) stream2go.Suscriber {
	return func(d interface {}){
		c<-d
	}
}
