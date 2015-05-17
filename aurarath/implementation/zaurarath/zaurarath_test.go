package zaurarath

import (
	"bytes"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/stream2go"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"testing"
	"time"
)

func TestPeerDiscover(t *testing.T) {

	iface1 := New([]byte("1111111111111111"))
	c1 := make(chan interface{})
	iface1.NewPeers().Listen(testlistener(c1))

	iface2 := New([]byte("2222222222222222"))
	c2 := make(chan interface{})
	iface2.NewPeers().Listen(testlistener(c2))

	data := (<-c1).(aurarath.Peer)
	if string(data.Id) != "2222222222222222" {
		t.Error("wrong id", string(data.Id))
	}
	data = (<-c2).(aurarath.Peer)
	iface1.Stop()
	iface2.Stop()

	if string(data.Id) != "1111111111111111" {
		t.Error("wrong id", string(data.Id))
	}

}

func TestHeartBeat(t *testing.T) {
	uuid1 := generateUuid()
	uuid2 := generateUuid()
	iface1 := New(uuid1)
	c1 := make(chan interface{})
	iface1.LeavingPeers().Where(isid(string(uuid2))).Listen(testlistener(c1))

	iface2 := New(uuid2)

	select {
	case <-time.After(10 * time.Second):
	case <-c1:
		t.Error("Peer leftto early")

	}
	iface2.Stop()

	iface1.Stop()

}
func TestPeerLeave(t *testing.T) {
	uuid1 := generateUuid()
	uuid2 := generateUuid()
	iface1 := New(uuid1)
	c1 := make(chan interface{})
	iface1.LeavingPeers().Where(isid(string(uuid2))).Listen(testlistener(c1))

	iface2 := New(uuid2)
	time.Sleep(2*time.Second)
	iface2.Stop()

	select {
	case <-time.After(10 * time.Second):
		t.Error("peer havent left")
	case <-c1:

	}
	iface1.Stop()

}

func TestPeerConnection(t *testing.T) {
	//	time.Sleep(1000*time.Millisecond)
	id1 := generateUuid()
	id2 := generateUuid()
	iface1 := New(id1)
	c1 := make(chan interface{})
	iface1.NewPeers().Where(isid(string(id2))).Listen(testlistener(c1))

	iface2 := New(id2)
	c2 := make(chan interface{})
	iface2.NewPeers().Where(isid(string(id1))).Listen(testlistener(c2))

	var home1 aurarath.Peer
	home1.Id = id1
	home1.Addresses = iface1.GetAdresses()
	var home2 aurarath.Peer
	home2.Id = id2
	home2.Addresses = iface2.GetAdresses()

	peer1 := (<-c1).(aurarath.Peer)
	if !bytes.Equal(peer1.Id, id2) {
		t.Error("wrong id", peer1.Id, id2)
	}
	peer2 := (<-c2).(aurarath.Peer)
	out1, err := iface1.Connect(home1, peer1.Addresses[0])
	if err != nil {
		t.Fatal(err)
	}
	out2, err := iface2.Connect(home2, peer2.Addresses[0])
	if err != nil {
		t.Fatal(err)
	}
	c3 := make(chan interface{})
	iface1.In().Listen(testlistener(c3))
	c4 := make(chan interface{})
	iface2.In().Listen(testlistener(c4))

	d1 := bytes.NewBufferString("Hello1")
	d2 := bytes.NewBufferString("Hello2")
	t1 := time.Now()
	out2.Add(aurarath.NewMessage(home2, 5, 6, []aurarath.Payload{aurarath.Payload{8, d2}}))
	m1 := (<-c3).(aurarath.Message)
	log.Println(time.Since(t1).Nanoseconds() / 1000)
	t2 := time.Now()
	out1.Add(aurarath.NewMessage(home1, 3, 4, []aurarath.Payload{aurarath.Payload{7, d1}}))
	m2 := (<-c4).(aurarath.Message)
	log.Println(time.Since(t2).Nanoseconds() / 1000)

	if string(m1.Payloads[0].Bytes.Bytes()) != "Hello2" {
		t.Error("wrong parameter")
	}
	if string(m2.Payloads[0].Bytes.Bytes()) != "Hello1" {
		t.Error("wrong parameter")
	}
	if m1.Protocol != 5 || m2.Protocol != 3 {
		t.Error("wrong protocoll")
	}
	if m1.Type != 6 || m2.Type != 4 {
		t.Error("wrong type")
	}

	if !bytes.Equal(m1.Sender.Id, home2.Id) || !bytes.Equal(m2.Sender.Id, home1.Id) {
		t.Error("wrong sender")
	}

	log.Println(m1)
	log.Println(m2)
}

func testlistener(c chan interface{}) stream2go.Suscriber {
	return func(d interface{}) {
		c <- d
	}
}

func eq(a interface{}, b interface{}) bool {
	e := map[interface{}]interface{}{}
	e[a] = nil
	_, f := e[b]
	return f
}

func isid(id string) stream2go.Filter {
	return func(d interface{}) bool {
		return string(d.(aurarath.Peer).Id) == id
	}
}

func BenchmarkTest(b *testing.B) {

	// RunParallel will create GOMAXPROCS goroutines
	// and distribute work among them.
	id1 := generateUuid()
	id2 := generateUuid()
	iface1 := New(id1)

	iface2 := New(id2)
	c2 := make(chan interface{})
	iface2.NewPeers().Where(isid(string(id1))).Listen(testlistener(c2))

	var home1 aurarath.Peer
	home1.Id = id1
	home1.Addresses = iface1.GetAdresses()
	var home2 aurarath.Peer
	home2.Id = id2
	home2.Addresses = iface2.GetAdresses()

	peer2 := (<-c2).(aurarath.Peer)

	out, _ := iface2.Connect(home2, peer2.Addresses[0])

	d := bytes.NewBufferString("Hello1")

	c3 := make(chan interface{})
	iface1.In().Listen(testlistener(c3))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out.Add(aurarath.NewMessage(home2, 5, 6, []aurarath.Payload{aurarath.Payload{8, d}}))
		<-c3
	}

}

func BenchmarkRound(b *testing.B) {
	id1 := generateUuid()
	id2 := generateUuid()
	iface1 := New(id1)
	c1 := make(chan interface{})
	iface1.NewPeers().Where(isid(string(id2))).Listen(testlistener(c1))

	iface2 := New(id2)
	c2 := make(chan interface{})
	iface2.NewPeers().Where(isid(string(id1))).Listen(testlistener(c2))

	var home1 aurarath.Peer
	home1.Id = id1
	home1.Addresses = iface1.GetAdresses()
	var home2 aurarath.Peer
	home2.Id = id2
	home2.Addresses = iface2.GetAdresses()

	peer1 := (<-c1).(aurarath.Peer)

	peer2 := (<-c2).(aurarath.Peer)
	out1, _ := iface1.Connect(home1, peer1.Addresses[0])

	out2, _ := iface2.Connect(home2, peer2.Addresses[0])

	c3 := make(chan interface{})
	iface1.In().Listen(testlistener(c3))
	c4 := make(chan interface{})
	iface2.In().Listen(testlistener(c4))

	d := bytes.NewBufferString("Hello1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out1.Add(aurarath.NewMessage(home2, 5, 6, []aurarath.Payload{aurarath.Payload{8, d}}))
		<-c4
		out2.Add(aurarath.NewMessage(home2, 5, 6, []aurarath.Payload{aurarath.Payload{8, d}}))
		<-c3
	}

}

func generateUuid() []byte {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return []byte{}
	}
	return Uuid[0:16]
}
