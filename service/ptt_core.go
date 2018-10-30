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
	"github.com/ailabstw/go-pttai/common"
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/pttdb"
)

func (p *BasePtt) GetVersion() (string, error) {
	return p.config.Version, nil
}

func (p *BasePtt) GetGitCommit() (string, error) {
	return p.config.GitCommit, nil
}

func (p *BasePtt) Shutdown() (bool, error) {
	p.notifyNodeStop.PassChan(struct{}{})
	return true, nil
}

func (p *BasePtt) Restart() (bool, error) {
	p.notifyNodeRestart.PassChan(struct{}{})
	return true, nil
}

/**********
 * Peer
 **********/

func (p *BasePtt) CountPeers() (*BackendCountPeers, error) {
	p.peerLock.RLock()
	defer p.peerLock.RUnlock()

	return &BackendCountPeers{
		MyPeers:        len(p.myPeers),
		ImportantPeers: len(p.importantPeers),
		MemberPeers:    len(p.memberPeers),
		RandomPeers:    len(p.randomPeers),
	}, nil
}

func (p *BasePtt) BEGetPeers() ([]*BackendPeer, error) {
	p.peerLock.RLock()
	defer p.peerLock.RUnlock()

	peerList := make([]*BackendPeer, 0, len(p.myPeers)+len(p.importantPeers)+len(p.memberPeers)+len(p.randomPeers))

	var backendPeer *BackendPeer
	for _, peer := range p.myPeers {
		backendPeer = PeerToBackendPeer(peer)
		peerList = append(peerList, backendPeer)
	}

	for _, peer := range p.importantPeers {
		backendPeer = PeerToBackendPeer(peer)
		peerList = append(peerList, backendPeer)
	}

	for _, peer := range p.memberPeers {
		backendPeer = PeerToBackendPeer(peer)
		peerList = append(peerList, backendPeer)
	}

	for _, peer := range p.randomPeers {
		backendPeer = PeerToBackendPeer(peer)
		peerList = append(peerList, backendPeer)
	}

	return peerList, nil
}

/**********
 * Entities
 **********/

func (p *BasePtt) CountEntities() (int, error) {
	return len(p.entities), nil
}

/**********
 * Join
 **********/

func (p *BasePtt) GetJoins() map[common.Address]*types.PttID {
	return p.joins
}

func (p *BasePtt) GetConfirmJoins() ([]*BackendConfirmJoin, error) {
	p.lockConfirmJoin.RLock()
	defer p.lockConfirmJoin.RUnlock()

	results := make([]*BackendConfirmJoin, len(p.confirmJoins))

	i := 0
	for _, confirmJoin := range p.confirmJoins {
		backendConfirmJoin := &BackendConfirmJoin{
			ID:         confirmJoin.JoinEntity.ID,
			Name:       confirmJoin.JoinEntity.Name,
			EntityID:   confirmJoin.Entity.GetID(),
			EntityName: []byte(confirmJoin.Entity.Name()),
			UpdateTS:   confirmJoin.UpdateTS,
			NodeID:     confirmJoin.Peer.GetID(),
			JoinType:   confirmJoin.JoinType,
		}
		results[i] = backendConfirmJoin

		i++
	}

	return results, nil
}

/**********
 * Op
 **********/

func (p *BasePtt) GetOps() map[common.Address]*types.PttID {
	return p.ops
}

/**********
 * MeOplog
 **********/

func (p *BasePtt) BEGetMeOplogList(logIDBytes []byte, limit int, listOrder pttdb.ListOrder) ([]*MeOplog, error) {

	var logID *types.PttID = nil
	if len(logIDBytes) != 0 {
		err := logID.Unmarshal(logIDBytes)
		if err != nil {
			return nil, err
		}
	}

	return p.GetMeOplogList(logID, limit, listOrder, types.StatusAlive)
}

func (p *BasePtt) BEGetPendingMeOplogMasterList(logIDBytes []byte, limit int, listOrder pttdb.ListOrder) ([]*MeOplog, error) {

	var logID *types.PttID = nil
	if len(logIDBytes) != 0 {
		err := logID.Unmarshal(logIDBytes)
		if err != nil {
			return nil, err
		}
	}

	return p.GetMeOplogList(logID, limit, listOrder, types.StatusPending)
}

func (p *BasePtt) BEGetPendingMeOplogInternalList(logIDBytes []byte, limit int, listOrder pttdb.ListOrder) ([]*MeOplog, error) {

	var logID *types.PttID = nil
	if len(logIDBytes) != 0 {
		err := logID.Unmarshal(logIDBytes)
		if err != nil {
			return nil, err
		}
	}

	return p.GetMeOplogList(logID, limit, listOrder, types.StatusInternalPending)
}

func (p *BasePtt) BEGetMeOplogMerkleNodeList(level MerkleTreeLevel, startKey []byte, limit int, listOrder pttdb.ListOrder) ([]*BackendMerkleNode, error) {

	merkleNodeList, err := p.GetMeOplogMerkleNodeList(level, startKey, limit, listOrder)
	if err != nil {
		return nil, err
	}

	results := make([]*BackendMerkleNode, len(merkleNodeList))
	for i, eachMerkleNode := range merkleNodeList {
		results[i] = MerkleNodeToBackendMerkleNode(eachMerkleNode)
	}

	return results, nil
}

/**********
 * MasterOplog
 **********/

func (p *BasePtt) BEGetMasterOplogList(logIDBytes []byte, limit int, listOrder pttdb.ListOrder) ([]*MasterOplog, error) {

	var logID *types.PttID = nil
	if len(logIDBytes) != 0 {
		err := logID.Unmarshal(logIDBytes)
		if err != nil {
			return nil, err
		}
	}

	return p.GetMasterOplogList(logID, limit, listOrder, types.StatusAlive)
}

func (p *BasePtt) BEGetMasterOplogMerkleNodeList(level MerkleTreeLevel, startKey []byte, limit int, listOrder pttdb.ListOrder) ([]*BackendMerkleNode, error) {

	merkleNodeList, err := p.GetMasterOplogMerkleNodeList(level, startKey, limit, listOrder)
	if err != nil {
		return nil, err
	}

	results := make([]*BackendMerkleNode, len(merkleNodeList))
	for i, eachMerkleNode := range merkleNodeList {
		results[i] = MerkleNodeToBackendMerkleNode(eachMerkleNode)
	}

	return results, nil
}

/**********
 * PttOplog
 **********/

func (p *BasePtt) BEGetPttOplogList(logIDBytes []byte, limit int, listOrder pttdb.ListOrder) ([]*PttOplog, error) {

	var logID *types.PttID = nil
	if len(logIDBytes) != 0 {
		err := logID.Unmarshal(logIDBytes)
		if err != nil {
			return nil, err
		}
	}

	return p.GetPttOplogList(logID, limit, listOrder, types.StatusAlive)
}

func (p *BasePtt) MarkPttOplogSeen() (types.Timestamp, error) {
	ts, err := types.GetTimestamp()
	if err != nil {
		return types.ZeroTimestamp, err
	}

	tsBytes, err := ts.Marshal()
	if err != nil {
		return types.ZeroTimestamp, err
	}

	err = dbMeta.Put(DBPttLogSeenPrefix, tsBytes)
	if err != nil {
		return types.ZeroTimestamp, err
	}

	return ts, nil
}

func (p *BasePtt) GetPttOplogSeen() (types.Timestamp, error) {
	tsBytes, err := dbMeta.Get(DBPttLogSeenPrefix)
	if err != nil {
		return types.ZeroTimestamp, nil
	}

	ts, err := types.UnmarshalTimestamp(tsBytes)
	if err != nil {
		return types.ZeroTimestamp, nil
	}

	return ts, nil
}
