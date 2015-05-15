package zaurarath

import (
	"sync"
	"time"
	"github.com/joernweissenborn/stream2go"
)

type PeerTracker struct {
	*sync.RWMutex
	trackedPeers map[string]time.Time
	deadPeers stream2go.StreamController
}

func NewPeerTracker() (pt *PeerTracker){
	pt = new(PeerTracker)
	pt.trackedPeers= []time.Time{}
	pt.deadPeers = stream2go.New()
	return
}

func (pt *PeerTracker) isTracked(d interface {}) (is bool) {
	pt.RLock()
	defer pt.Unlock()
	var s Signal
	s, is = d.(Signal)
	if !is {
		return
	}
	_,is := pt.trackedPeers[string(s.Data[:16])]
	return
}


func (pt *PeerTracker) TrackPeer(Id string) {

	pt.Lock()
	defer  pt.Unlock()

	pt.trackedPeers[Id] = time.Now()

}


func (pt *PeerTracker) Track() {

	var checktime = 1 * time.Second

	t := time.NewTimer(checktime)
	for {
		select {

		case <-t.C:
			t.Reset(checktime)
			pt.checkTimeout()
		}
	}

}



func (pt *PeerTracker) checkTimeout() {
	defer pt.Unlock()
	for id, p := range pt.trackedPeers {
		pt.RLock()
		lastCheck := time.Since(p).Seconds()
		pt.Unlock()
		if lastCheck < 5.0 {
			pt.PeerDead(id)
		}
	}
	return
}

func (pt *PeerTracker) PeerDead(Id string) {
	pt.Lock()
	defer pt.Unlock()
	delete(pt.trackedPeers,Id)
}
func (pt *PeerTracker) Heartbeat(d interface {}) {
	pt.WLock()
	defer pt.Unlock()
	var s Signal
	s = d.(Signal)
	pt.trackedPeers[string(s.Data[:16])] = time.Now()

	return
}
