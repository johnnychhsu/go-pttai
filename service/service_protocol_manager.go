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

	"github.com/ailabstw/go-pttai/common/types"
)

/*
ServiceProtocolManager manage service-level operations.

ServiceProtocolManager includes peers-of-services and the corresponding entities.
both are dynamically allocated / deallocated.

When there is a new peer: have all the existing entities trying to register the peer.
When a peer disappear: have all the existing entities trying to unregister the peer.

When there is a new entity: trying to register all the peers.
When a peer disappear: trying to unregister all the peers.
*/
type ServiceProtocolManager interface {
	Start() error
	Stop() error

	// entities
	Entities() map[types.PttID]Entity
	Entity(id *types.PttID) Entity

	RegisterEntity(id *types.PttID, e Entity) error
	UnregisterEntity(id *types.PttID) error

	Ptt() Ptt
	Service() Service
}

type BaseServiceProtocolManager struct {
	lock     sync.RWMutex
	entities map[types.PttID]Entity

	noMorePeers chan struct{}

	ptt     Ptt
	service Service
}

func NewBaseServiceProtocolManager(ptt Ptt, service Service) (*BaseServiceProtocolManager, error) {
	spm := &BaseServiceProtocolManager{
		entities: make(map[types.PttID]Entity),

		noMorePeers: ptt.NoMorePeers(),
		ptt:         ptt,

		service: service,
	}

	return spm, nil
}

func (spm *BaseServiceProtocolManager) Start() error {
	return nil
}

func (spm *BaseServiceProtocolManager) Stop() error {
	return nil
}

func (spm *BaseServiceProtocolManager) Entities() map[types.PttID]Entity {
	return spm.entities
}

func (spm *BaseServiceProtocolManager) Ptt() Ptt {
	return spm.ptt
}

func (spm *BaseServiceProtocolManager) Service() Service {
	return spm.service
}

func (spm *BaseServiceProtocolManager) Entity(id *types.PttID) Entity {
	spm.lock.RLock()
	defer spm.lock.RUnlock()

	entity, ok := spm.entities[*id]
	if !ok {
		return nil
	}
	return entity
}

/*
RegisterEntity register the entity to the service

need to do lock in the beginning because need to update entitiesByPeerID
*/
func (spm *BaseServiceProtocolManager) RegisterEntity(id *types.PttID, e Entity) error {
	spm.lock.Lock()
	defer spm.lock.Unlock()

	_, ok := spm.entities[*id]
	if ok {
		return ErrEntityAlreadyRegistered
	}

	spm.entities[*id] = e
	e.PM().SetNoMorePeers(spm.noMorePeers)

	return nil
}

func (spm *BaseServiceProtocolManager) UnregisterEntity(id *types.PttID) error {
	spm.lock.Lock()
	defer spm.lock.Unlock()

	_, ok := spm.entities[*id]
	if !ok {
		return ErrEntityNotRegistered
	}

	delete(spm.entities, *id)

	return nil
}
