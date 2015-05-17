package zaurarath

import (
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/stream2go"
	"sync"
	"time"
)

type PeerTracker struct {
	*sync.RWMutex
	trackedPeers map[string]TrackedPeer
	deadPeers    stream2go.StreamController
}

func NewPeerTracker() (pt *PeerTracker) {
	pt = new(PeerTracker)
	pt.RWMutex = new(sync.RWMutex)
	pt.trackedPeers = map[string]TrackedPeer{}
	pt.deadPeers = stream2go.New()
	return
}

func (pt *PeerTracker) isKnown(d interface{}) (is bool) {
	pt.RLock()
	defer pt.RUnlock()
	peer := d.(aurarath.Peer)

	_, is = pt.trackedPeers[string(peer.Id)]

	return
}

func (pt *PeerTracker) isKnownAddr(d interface{}) (is bool, add aurarath.Address) {
	pt.RLock()
	defer pt.RUnlock()
	peer := d.(aurarath.Peer)

	p, is := pt.trackedPeers[string(peer.Id)]
	if is {
		add = p.Peer.Addresses[0]
	}
	return
}

func (pt *PeerTracker) add(d interface{}) {
	pt.Lock()
	defer pt.Unlock()
	peer := d.(aurarath.Peer)

	var p TrackedPeer
	p.Peer = peer
	pt.trackedPeers[string(p.Id)] = p

	return
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

	for _, p := range pt.trackedPeers {
		pt.RLock()
		lastCheck := time.Since(p.lastCheckin).Seconds()
		pt.RUnlock()

		if lastCheck > 5.0 {
			pt.PeerDead(p.Peer)
		}
	}
	return
}

func (pt *PeerTracker) PeerDead(p aurarath.Peer) {
	pt.Lock()
	defer pt.Unlock()
	delete(pt.trackedPeers, string(p.Id))
	pt.deadPeers.Add(p)
}
func (pt *PeerTracker) Heartbeat(d interface{}) {
	pt.Lock()
	defer pt.Unlock()
	peer := d.(aurarath.Peer)

	p := pt.trackedPeers[string(peer.Id)]
	p.lastCheckin = time.Now()
	pt.trackedPeers[string(peer.Id)] = p
	return
}

type TrackedPeer struct {
	aurarath.Peer
	lastCheckin time.Time
}
