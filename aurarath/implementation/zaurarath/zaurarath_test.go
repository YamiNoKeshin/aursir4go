package zaurarath

import (
	"testing"
	"github.com/joernweissenborn/stream2go"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"bytes"
	"log"
	"time"
)


func TestPeerDiscover(t *testing.T){

	iface1 := New([]byte("1111111111111111"))
	c1 := make(chan interface {})
	iface1.NewPeers().Listen(testlistener(c1))

	iface2 := New([]byte("2222222222222222"))
	c2 := make(chan interface {})
	iface2.NewPeers().Listen(testlistener(c2))

	data := (<-c1).(aurarath.Peer)
	if string(data.Id) != "2222222222222222" {
		t.Error("wrong id",string(data.Id))
	}
	data = (<-c2).(aurarath.Peer)
	iface1.Stop()
	 iface2.Stop()

	if string(data.Id) != "1111111111111111" {
		t.Error("wrong id",string(data.Id))
	}

}


func TestPeerLeave(t *testing.T){

	iface1 := New([]byte("1111111111111111"))
	c1 := make(chan interface {})
	iface1.LeavingPeers().Where(isid("2222222222222222")).Listen(testlistener(c1))

	iface2 := New([]byte("2222222222222222"))

	 iface2.Stop()

	select {
	case <- time.After(51*time.Second):
		t.Error("Beaon didnt stop")
	case <- c1:

	}
	iface1.Stop()

}

func TestPeerConnection(t *testing.T){
//	time.Sleep(1000*time.Millisecond)
	iface1 := New([]byte("3111111111111111"))
	c1 := make(chan interface {})
	iface1.NewPeers().Where(isid("4222222222222222")).Listen(testlistener(c1))

	iface2 := New([]byte("4222222222222222"))
	c2 := make(chan interface {})
	iface2.NewPeers().Where(isid("3111111111111111")).Listen(testlistener(c2))

	var home1 aurarath.Peer
	home1.Id = []byte("3111111111111111")
	home1.Addresses = iface1.GetAdresses()
	var home2 aurarath.Peer
	home2.Id = []byte("4222222222222222")
	home2.Addresses = iface2.GetAdresses()

	peer1 := (<-c1).(aurarath.Peer)
	if string(peer1.Id) != "4222222222222222" {
		t.Error("wrong id",string(peer1.Id))
	}
	peer2 := (<-c2).(aurarath.Peer)
	out1, err := iface1.Connect(home1,peer1.Addresses[0])
	if err != nil {
		t.Fatal(err)
	}
	out2, err := iface2.Connect(home2,peer2.Addresses[0])
	if err != nil {
		t.Fatal(err)
	}
	c3 := make(chan interface {})
	iface1.In().Listen(testlistener(c3))
	c4 := make(chan interface {})
	iface2.In().Listen(testlistener(c4))


	d1 := bytes.NewBufferString("Hello1")
	d2 := bytes.NewBufferString("Hello2")
	t1:= time.Now()
	out2.Add(aurarath.NewMessage(home2, 5,6, []aurarath.Payload{aurarath.Payload{8,d2}}))
	m1 := (<-c3).(aurarath.Message)
	log.Println(time.Since(t1).Nanoseconds()/1000)
	t2:= time.Now()
	out1.Add(aurarath.NewMessage(home1, 3,4, []aurarath.Payload{aurarath.Payload{7,d1}}))
	m2 := (<-c4).(aurarath.Message)
	log.Println(time.Since(t2).Nanoseconds()/1000)


	if string(m1.Payloads[0].Bytes.Bytes()) != "Hello2" {
		t.Error("wrong parameter")
	}
	if string(m2.Payloads[0].Bytes.Bytes()) != "Hello1" {
		t.Error("wrong parameter")
	}
	if m1.Protocol != 5 || m2.Protocol != 3 {
		t.Error("wrong protocoll")
	}
	if m1.Type != 6 || m2.Type != 4{
		t.Error("wrong type")
	}
	if !bytes.Equal(m1.Sender.Id ,home2.Id) || !bytes.Equal(m2.Sender.Id ,home1.Id){
		t.Error("wrong sender")
	}

	log.Println(m1)
	log.Println(m2)
}

func testlistener(c chan interface {}) stream2go.Suscriber {
	return func(d interface {}){
		c<-d
	}
}

func isid(id string) stream2go.Filter {
	return func(d interface {})bool{
		return string(d.(aurarath.Peer).Id) == id
	}
}


