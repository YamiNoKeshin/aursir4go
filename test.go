package main

import (
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"bytes"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/stream2go"
	"github.com/joernweissenborn/aursir4go/aurarath/implementation/zaurarath"
	"sync"
)

type bla struct {
	h string
}

func main() {
	wg := new(sync.WaitGroup)
	fails := 0
	for i := 0; i<100 ;i++ {
		wg.Add(1)
//		time.Sleep(100*time.Millisecond)
		go stresstest(wg)
		id := generateUuid()
		if bytes.Equal(id[:1],[]byte{0}) {
			fails++
		}
	}
	log.Println("ID FAIL", fails)
	wg.Wait()
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

func eq(a interface{}, b interface{}) bool {
	e := map[interface{}]interface{}{}
	e[a] = nil
	_, f := e[b]
	return f
}

func stresstest(wg *sync.WaitGroup) {
	id1 := generateUuid()
	id2 := generateUuid()
	iface1 := zaurarath.New(id1)
	c1 := make(chan interface{})
	iface1.NewPeers().Where(isid(string(id2))).Listen(testlistener(c1))

	iface2 := zaurarath.New(id2)
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

	out1, err := iface1.Connect(home1, peer1.Addresses[0])
	if err != nil {
	}
	defer out1.Close()
	out2, err := iface2.Connect(home2, peer2.Addresses[0])
	if err != nil {
	}
	defer out2.Close()
	c3 := make(chan interface{})
	iface1.In().Listen(testlistener(c3))
	c4 := make(chan interface{})
	iface2.In().Listen(testlistener(c4))

	d1 := bytes.NewBufferString("Hello1")
	d2 := bytes.NewBufferString("Hello2")
	out2.Add(aurarath.NewMessage(home2, 5, 6, []aurarath.Payload{aurarath.Payload{8, d2}}))
	m1 := (<-c3).(aurarath.Message)

	out1.Add(aurarath.NewMessage(home1, 3, 4, []aurarath.Payload{aurarath.Payload{7, d1}}))
	m2 := (<-c4).(aurarath.Message)

	if string(m1.Payloads[0].Bytes.Bytes()) != "Hello2" {
		log.Println("wrong parameter")
	}
	if string(m2.Payloads[0].Bytes.Bytes()) != "Hello1" {
		log.Println("wrong parameter")
	}
	if m1.Protocol != 5 || m2.Protocol != 3 {
		log.Println("wrong protocoll")
	}
	if m1.Type != 6 || m2.Type != 4 {
		log.Println("wrong type")
	}
	if !bytes.Equal(peer1.Id, id2) {
		log.Println("wrong id", peer1.Id, id2)
	}

	if !bytes.Equal(m1.Sender.Id, home2.Id) || !bytes.Equal(m2.Sender.Id, home1.Id) {
		log.Println("wrong sender")
		log.Println(string(m1.Sender.Id))
		log.Println(m1.Sender.Id)
		log.Println(len(m1.Sender.Id))
		log.Println(string(home2.Id))
		log.Println(home2.Id)
		log.Println(len(home2.Id))
		log.Println(string(m2.Sender.Id))
		log.Println(m2.Sender.Id)
		log.Println(len(m2.Sender.Id))
		log.Println(string(home1.Id))
		log.Println(home1.Id)
		log.Println(len(home1.Id))
	}

	if !bytes.Equal(peer2.Id, id1) {
		log.Println("wrong id", peer2.Id, id1)
	}

	out2.Add(aurarath.NewMessage(home2, 5, 6, []aurarath.Payload{aurarath.Payload{8, d2}}))
	(<-c3)
	out1.Add(aurarath.NewMessage(home1, 3, 4, []aurarath.Payload{aurarath.Payload{7, d1}}))
	(<-c4)
	wg.Done()
}



func isid(id string) stream2go.Filter {
	return func(d interface{}) bool {
		return string(d.(aurarath.Peer).Id) == id
	}
}

func testlistener(c chan interface{}) stream2go.Suscriber {
	return func(d interface{}) {
		c <- d
	}
}
