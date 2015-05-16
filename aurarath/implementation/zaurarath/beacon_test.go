package zaurarath

import (
	"github.com/joernweissenborn/future2go"
	"testing"
)


func TetBeacon(t *testing.T){
	b1 := NewBeacon([]byte("1234"))
	defer b1.Stop()
	c1 := make(chan interface {})
	b1.Signals().First().Then(testcompleter(c1))
	b2 := NewBeacon([]byte("HALLO"))
	defer b2.Stop()
	c2 := make(chan interface {})
	b2.Signals().First().Then(testcompleter(c2))
	b1.Run()
	b2.Run()
	data := string((<-c1).(Signal).Data)
	if data != "HALLO" {
		t.Error("got wrong data, needed 'HALLO', got",data)
	}
	data = string((<-c2).(Signal).Data)

	if string((<-c2).(Signal).Data) != "1234" {
		t.Error("got wrong data",data)
	}
}


func testcompleter(c chan interface {}) future2go.CompletionFunc {
	return func(d interface {})interface {}{
		c<-d
		return nil
	}
}
