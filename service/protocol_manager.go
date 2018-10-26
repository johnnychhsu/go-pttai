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
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/event"
	"github.com/ailabstw/go-pttai/pttdb"
)

type ProtocolManager interface {
	Start() error
	Stop() error

	// implemented in base-protocol-manager
	// event-mux
	EventMux() *event.TypeMux

	// master
	Master0Hash() []byte
	IsMaster(id *types.PttID) bool
	SetNewestMasterLogID(id *types.PttID)
	GetNewestMasterLogID() *types.PttID

	// entity
	Entity() Entity

	// ptt
	Ptt() Ptt

	// db
	DB() *pttdb.LDBBatch

	// is-start
	IsStart() bool
}

type BaseProtocolManager struct {
	// event-mux
	eventMux *event.TypeMux

	// master-log-id
	newestMasterLogID *types.PttID

	// entity
	entity Entity

	// ptt
	ptt Ptt

	// db
	db *pttdb.LDBBatch

	// isStart
	isStart bool
}

func NewBaseProtocolManager(ptt Ptt, renewOpKeySeconds uint64, expireOpKeySeconds uint64, maxSyncRandomSeconds int, minSyncRandomSeconds int, e Entity, db *pttdb.LDBBatch) (*BaseProtocolManager, error) {

	pm := &BaseProtocolManager{
		// event-mux
		eventMux: new(event.TypeMux),

		// entity
		entity: e,

		// ptt
		ptt: ptt,

		// db
		db: db,
	}

	// master-log-id
	newestMasterLogID, err := pm.loadNewestMasterLogID()
	if err != nil {
		return nil, err
	}

	pm.newestMasterLogID = newestMasterLogID

	return pm, nil

}

func (pm *BaseProtocolManager) Start() error {
	return nil
}

func (pm *BaseProtocolManager) Stop() error {
	pm.eventMux.Stop()

	return nil
}

func (pm *BaseProtocolManager) EventMux() *event.TypeMux {
	return pm.eventMux
}

func (pm *BaseProtocolManager) Ptt() Ptt {
	return pm.ptt
}

func (pm *BaseProtocolManager) DB() *pttdb.LDBBatch {
	return pm.db
}

func (pm *BaseProtocolManager) Entity() Entity {
	return pm.entity
}

func (pm *BaseProtocolManager) IsStart() bool {
	return pm.isStart
}
