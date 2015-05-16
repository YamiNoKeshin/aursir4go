package zaurarath

import (
	"github.com/joernweissenborn/future2go"
	"testing"
)


func TestBeacon(t *testing.T) {
	b1 := NewBeacon([]byte("1234"), 9999)
	defer b1.Stop()
	c1 := make(chan interface{})
	b1.Signals().First().Then(testcompleter(c1))
	b2 := NewBeacon([]byte("HALLO"), 9999)
	defer b2.Stop()
	c2 := make(chan interface{})
	b2.Signals().First().Then(testcompleter(c2))
	b1.Run()
	b2.Run()
	data := (<-c1).(Signal)
	if string(data.Data) != "HALLO" {
		t.Error("got wrong data, needed 'HALLO', got", data)
	}
	data2 := (<-c2).(Signal)

	if string(data2.Data) != "1234" {
		t.Error("got wrong data", data)
	}
}


func testcompleter(c chan interface {}) future2go.CompletionFunc {
	return func(d interface {})interface {}{
		c<-d
		return nil
	}
}


func checkip(c chan interface {}) future2go.CompletionFunc {
	return func(d interface {})interface {}{
		c<-d
		return nil
	}
}
