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
	"time"

	"github.com/ailabstw/go-pttai/common"
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/event"
	"github.com/ailabstw/go-pttai/p2p/discover"
	"github.com/ailabstw/go-pttai/pttdb"
)

type ProtocolManager interface {
	Start() error
	Stop() error

	HandleMessage(op OpType, dataBytes []byte, peer *PttPeer) error

	Sync(peer *PttPeer) error

	// implemented in base-protocol-manager
	// event-mux
	EventMux() *event.TypeMux

	// peers
	Peers() *PttPeerSet
	PeerList() []*PttPeer

	NewPeerCh() chan *PttPeer
	NoMorePeers() chan struct{}
	SetNoMorePeers(noMorePeers chan struct{})

	RegisterPeer(peer *PttPeer) error
	UnregisterPeer(peer *PttPeer) error

	CountPeers() (int, error)

	IsMyDevice(peer *PttPeer) bool
	IsImportantPeer(peer *PttPeer) bool
	IsMemberPeer(peer *PttPeer) bool

	IsFitPeer(peer *PttPeer) PeerType

	IsSuspiciousID(id *types.PttID, nodeID *discover.NodeID) bool

	IsGoodID(id *types.PttID, nodeID *discover.NodeID) bool

	// sync
	QuitSync() chan struct{}
	SetQuitSync(quitSync chan struct{})

	SyncWG() *sync.WaitGroup

	ForceSyncCycle() time.Duration
	SetForceSyncCycle()

	// join
	JoinKeyInfos() []*KeyInfo
	GetJoinKeyInfo(hash *common.Address) (*KeyInfo, error)
	GetJoinKey() (*KeyInfo, error)

	GenerateJoinKeyInfoLoop() error

	ApproveJoin(joinEntity *JoinEntity, keyInfo *KeyInfo, peer *PttPeer) (*KeyInfo, interface{}, error)

	GetJoinType(hash *common.Address) (JoinType, error)

	// op-key
	GetOpKeyInfo(hash *common.Address) (*KeyInfo, error)
	GetOpKey() (*KeyInfo, error)
	GetNewestOpKey() (*KeyInfo, error)
	RegisterOpKeyInfo(keyInfo *KeyInfo) error

	DBOpKeyLock() *types.LockMap

	OpKeyInfos() []*KeyInfo

	SaveOpKeyInfo(opKeyInfo *KeyInfo) error

	RevokeOpKeyInfo(revokeOpKeyInfo *KeyInfo) error

	RenewOpKeySeconds() uint64
	ExpireOpKeySeconds() uint64

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

	// peers
	peers *PttPeerSet

	newPeerCh   chan *PttPeer
	noMorePeers chan struct{}

	// sync
	quitSync chan struct{}
	syncWG   *sync.WaitGroup

	maxSyncRandomSeconds int
	minSyncRandomSeconds int

	forceSyncCycle time.Duration

	// join
	// we may use map to speed up the lookup from hash, but map may significantly increasing the mem. Considering that the keys should be kept < 10, we use sorted-by-update-ts array.

	lockJoinKeyInfo sync.RWMutex
	joinKeyInfos    []*KeyInfo

	// op-key
	lockOpKeyInfo sync.RWMutex
	opKeyInfos    []*KeyInfo

	opKeyChan     chan *KeyInfo
	revokeKeyChan chan *KeyInfo

	renewOpKeySeconds  uint64
	expireOpKeySeconds uint64

	dbOpKeyLock *types.LockMap

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

	peerSet, err := NewPttPeerSet()
	if err != nil {
		return nil, err
	}

	dbOpKeyLock, err := types.NewLockMap(SleepTimeOpKeyLock)
	if err != nil {
		return nil, err
	}

	pm := &BaseProtocolManager{
		// event-mux
		eventMux: new(event.TypeMux),

		// peers
		peers: peerSet,

		// new-peer-ch
		newPeerCh: make(chan *PttPeer),

		// sync
		maxSyncRandomSeconds: maxSyncRandomSeconds,
		minSyncRandomSeconds: minSyncRandomSeconds,

		syncWG: ptt.SyncWG(),

		// join-key
		joinKeyInfos: make([]*KeyInfo, 0),

		// op-key
		renewOpKeySeconds:  renewOpKeySeconds,
		expireOpKeySeconds: expireOpKeySeconds,

		dbOpKeyLock: dbOpKeyLock,

		// entity
		entity: e,

		// ptt
		ptt: ptt,

		// db
		db: db,
	}

	// op-key
	opKeyInfos, err := pm.loadOpKeyInfos(expireOpKeySeconds)
	if err != nil {
		return nil, err
	}

	pm.opKeyInfos = opKeyInfos

	// master-log-id
	newestMasterLogID, err := pm.loadNewestMasterLogID()
	if err != nil {
		return nil, err
	}

	pm.newestMasterLogID = newestMasterLogID

	return pm, nil

}

func (pm *BaseProtocolManager) Start() error {
	pm.isStart = true

	return nil
}

func (pm *BaseProtocolManager) Stop() error {
	pm.eventMux.Stop()

	return nil
}

func (pm *BaseProtocolManager) HandleMessage(op OpType, dataBytes []byte, peer *PttPeer) error {
	return nil
}

func (pm *BaseProtocolManager) Sync(peer *PttPeer) error {
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

func (pm *BaseProtocolManager) OpKeyInfos() []*KeyInfo {
	return pm.opKeyInfos
}

func (pm *BaseProtocolManager) SaveOpKeyInfo(opKeyInfo *KeyInfo) error {
	return opKeyInfo.Save(pm.db)
}

func (pm *BaseProtocolManager) Entity() Entity {
	return pm.entity
}

func (pm *BaseProtocolManager) GetJoinType(hash *common.Address) (JoinType, error) {
	return JoinTypeInvalid, ErrInvalidData
}

func (pm *BaseProtocolManager) DBOpKeyLock() *types.LockMap {
	return pm.dbOpKeyLock
}
