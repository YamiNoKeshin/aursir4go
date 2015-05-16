package zaurarath

import (
	"sync"
	"time"
	"github.com/joernweissenborn/stream2go"
)

type PeerTracker struct {
	*sync.RWMutex
	trackedPeers map[string]TrackedPeer
	deadPeers stream2go.StreamController
}

func NewPeerTracker() (pt *PeerTracker){
	pt = new(PeerTracker)
	pt.RWMutex = new(sync.RWMutex)
	pt.trackedPeers= map[string]TrackedPeer{}
	pt.deadPeers = stream2go.New()
	return
}

func (pt *PeerTracker) isTracked(d interface {}) (is bool) {
	pt.RLock()
	defer pt.RUnlock()
	var s Signal
	s, is = d.(Signal)
	if !is {
		return
	}
	p,is := pt.trackedPeers[string(s.Data[:16])]
	if is {
		is = p.track
	}
	return
}

func (pt *PeerTracker) isKnown(d interface {}) (is bool) {
	pt.RLock()
	defer pt.RUnlock()
	var s Signal
	s, is = d.(Signal)
	if !is {
		return
	}
	_,is = pt.trackedPeers[string(s.Data[:16])]

	return
}

func (pt *PeerTracker) add(d interface {}) {
	pt.Lock()
	defer pt.Unlock()
	s, is := d.(Signal)
	if !is {
		return
	}
	var p TrackedPeer
	pt.trackedPeers[string(s.Data[:16])] = p

	return
}


func (pt *PeerTracker) TrackPeer(Id string) {

	pt.Lock()
	defer  pt.Unlock()

	p, f := pt.trackedPeers[Id]
	if f {
		p.track = true
		p.lastCheckin = time.Now()
	}

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
	for id, p := range pt.trackedPeers {
		pt.RLock()
		if !p.track {continue}
		lastCheck := time.Since(p.lastCheckin).Seconds()
		pt.RUnlock()
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
	pt.deadPeers.Add(Id)
}
func (pt *PeerTracker) Heartbeat(d interface {}) {
	pt.Lock()
	defer pt.Unlock()
	var s Signal
	s = d.(Signal)

	p := pt.trackedPeers[string(s.Data[:16])]
	p.lastCheckin = time.Now()
	return
}

type TrackedPeer struct {
	track bool
	lastCheckin time.Time
}
