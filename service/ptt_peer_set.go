// Copyright 2018 The go-pttai Authors
// This file is part of the go-pttai library.
//
// The go-pttai library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-pttai library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-pttai library. If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"sync"

	"github.com/ailabstw/go-pttai/p2p/discover"
)

type PttPeerSet struct {
	peerTypes map[discover.NodeID]PeerType

	mePeers           map[discover.NodeID]*PttPeer
	mePeerList        []*PttPeer
	importantPeers    map[discover.NodeID]*PttPeer
	importantPeerList []*PttPeer
	memberPeers       map[discover.NodeID]*PttPeer
	memberPeerList    []*PttPeer
	peerList          []*PttPeer

	lock   sync.RWMutex
	closed bool
}

func NewPttPeerSet() (*PttPeerSet, error) {
	return &PttPeerSet{
		peerTypes: make(map[discover.NodeID]PeerType),

		mePeers:           make(map[discover.NodeID]*PttPeer),
		importantPeers:    make(map[discover.NodeID]*PttPeer),
		memberPeers:       make(map[discover.NodeID]*PttPeer),
		mePeerList:        make([]*PttPeer, 0),
		importantPeerList: make([]*PttPeer, 0),
		memberPeerList:    make([]*PttPeer, 0),
		peerList:          make([]*PttPeer, 0),
	}, nil
}

func (ps *PttPeerSet) MePeers(isLocked bool) map[discover.NodeID]*PttPeer {
	if !isLocked {
		ps.RLock()
		defer ps.RUnlock()
	}
	return ps.mePeers
}

func (ps *PttPeerSet) ImportantPeers() map[discover.NodeID]*PttPeer {
	return ps.importantPeers
}

func (ps *PttPeerSet) MemberPeers() map[discover.NodeID]*PttPeer {
	return ps.memberPeers
}

func (ps *PttPeerSet) Lock() {
	ps.lock.Lock()
}

func (ps *PttPeerSet) Unlock() {
	ps.lock.Unlock()
}

func (ps *PttPeerSet) RLock() {
	ps.lock.RLock()
}

func (ps *PttPeerSet) RUnlock() {
	ps.lock.RUnlock()
}

func (ps *PttPeerSet) IsClosed() bool {
	return ps.closed
}

func (ps *PttPeerSet) Register(peer *PttPeer, peerType PeerType, isLocked bool) error {
	if !isLocked {
		ps.Lock()
		defer ps.Unlock()
	}

	if ps.closed {
		return ErrClosed
	}

	id := peer.ID()
	pid := peer.GetID()
	origPeerType, ok := ps.peerTypes[id]
	if ok {
		origPeer := ps.GetPeerWithPeerType(pid, origPeerType, true)
		if origPeer != nil && origPeer != peer {
			return ErrAlreadyRegistered
		}

		if origPeerType == peerType {
			return nil
		}
	}

	ps.peerTypes[id] = peerType

	switch origPeerType {
	case PeerTypeMe:
		delete(ps.mePeers, id)
	case PeerTypeImportant:
		delete(ps.importantPeers, id)
	case PeerTypeMember:
		delete(ps.memberPeers, id)
	}

	switch peerType {
	case PeerTypeMe:
		ps.mePeers[id] = peer
	case PeerTypeImportant:
		ps.importantPeers[id] = peer
	case PeerTypeMember:
		ps.memberPeers[id] = peer
	}

	if origPeerType == PeerTypeMe || peerType == PeerTypeMe {
		ps.mePeerList = ps.PeersToPeerList(ps.mePeers, true)
	}
	if origPeerType == PeerTypeImportant || peerType == PeerTypeImportant {
		ps.importantPeerList = ps.PeersToPeerList(ps.importantPeers, true)
	}
	if origPeerType == PeerTypeMember || peerType == PeerTypeMember {
		ps.memberPeerList = ps.PeersToPeerList(ps.memberPeers, true)
	}

	ps.peerList = append(append(ps.mePeerList, ps.importantPeerList...), ps.memberPeerList...)

	return nil
}

func (ps *PttPeerSet) GetPeerWithPeerType(id *discover.NodeID, peerType PeerType, isLocked bool) *PttPeer {
	if !isLocked {
		ps.RLock()
		defer ps.RUnlock()
	}

	switch peerType {
	case PeerTypeMe:
		return ps.mePeers[*id]
	case PeerTypeImportant:
		return ps.importantPeers[*id]
	case PeerTypeMember:
		return ps.memberPeers[*id]
	}

	return nil
}

func (ps *PttPeerSet) Unregister(peer *PttPeer, isLocked bool) error {
	if !isLocked {
		ps.Lock()
		defer ps.Unlock()
	}

	id := peer.ID()
	origPeerType, ok := ps.peerTypes[id]
	if !ok {
		return ErrNotRegistered
	}

	switch origPeerType {
	case PeerTypeMe:
		delete(ps.mePeers, id)
	case PeerTypeImportant:
		delete(ps.importantPeers, id)
	case PeerTypeMember:
		delete(ps.memberPeers, id)
	}

	delete(ps.peerTypes, id)

	if origPeerType == PeerTypeMe {
		ps.mePeerList = ps.PeersToPeerList(ps.mePeers, true)
	}
	if origPeerType == PeerTypeImportant {
		ps.importantPeerList = ps.PeersToPeerList(ps.importantPeers, true)
	}
	if origPeerType == PeerTypeMember {
		ps.memberPeerList = ps.PeersToPeerList(ps.memberPeers, true)
	}

	ps.peerList = append(append(ps.mePeerList, ps.importantPeerList...), ps.memberPeerList...)

	return nil
}

func (ps *PttPeerSet) PeersToPeerList(peers map[discover.NodeID]*PttPeer, isLocked bool) []*PttPeer {
	if !isLocked {
		ps.lock.Lock()
		defer ps.lock.Unlock()
	}

	lenPeers := len(peers)

	peerList := make([]*PttPeer, 0, lenPeers)
	i := 0
	for _, peer := range peers {
		peerList[i] = peer
		i++
	}

	return peerList
}

func (ps *PttPeerSet) MePeerList(isLocked bool) []*PttPeer {
	if !isLocked {
		ps.RLock()
		defer ps.RUnlock()
	}

	return ps.mePeerList
}

func (ps *PttPeerSet) ImportantPeerList(isLocked bool) []*PttPeer {
	if !isLocked {
		ps.RLock()
		defer ps.RUnlock()
	}

	return ps.importantPeerList
}

func (ps *PttPeerSet) MemberPeerList(isLocked bool) []*PttPeer {
	if !isLocked {
		ps.RLock()
		defer ps.RUnlock()
	}

	return ps.memberPeerList
}

func (ps *PttPeerSet) PeerList(isLocked bool) []*PttPeer {
	if !isLocked {
		ps.RLock()
		defer ps.RUnlock()
	}

	return ps.peerList
}

// Peer retrieves the registered peer with the given id.
func (ps *PttPeerSet) Peer(id *discover.NodeID, isLocked bool) *PttPeer {
	if !isLocked {
		ps.RLock()
		defer ps.RUnlock()
	}

	peerType, ok := ps.peerTypes[*id]
	if !ok {
		return nil
	}

	switch peerType {
	case PeerTypeMe:
		return ps.mePeers[*id]
	case PeerTypeImportant:
		return ps.importantPeers[*id]
	case PeerTypeMember:
		return ps.memberPeers[*id]
	}

	return nil
}

// Len returns if the current number of peers in the set.
func (ps *PttPeerSet) Len(isLocked bool) int {
	if !isLocked {
		ps.RLock()
		defer ps.RUnlock()
	}

	return len(ps.peerTypes)
}

// Close disconnects all peers.
// No new peers can be registered after Close has returned.
func (ps *PttPeerSet) Close(isLocked bool) {
	if !isLocked {
		ps.Lock()
		defer ps.Unlock()
	}

	ps.closed = true

	ps.peerTypes = make(map[discover.NodeID]PeerType)
	ps.mePeers = make(map[discover.NodeID]*PttPeer)
	ps.importantPeers = make(map[discover.NodeID]*PttPeer)
	ps.memberPeers = make(map[discover.NodeID]*PttPeer)

}
