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
	"encoding/json"

	"github.com/ailabstw/go-pttai/common/types"
)

type SyncOplogAck struct {
	TS    types.Timestamp
	Nodes []*MerkleNode `json:"N"`
}

func (pm *BaseProtocolManager) SyncOplogAck(lastSyncTime types.Timestamp, merkle *Merkle, op OpType, peer *PttPeer) error {
	now, err := types.GetTimestamp()
	if err != nil {
		return err
	}

	offsetHourTS, _ := lastSyncTime.ToHRTimestamp()
	nodes, err := merkle.GetMerkleTreeListByLevel(MerkleTreeLevelNow, offsetHourTS, now)
	if err != nil {
		return err
	}

	syncOplogAck := &SyncOplogAck{
		TS:    offsetHourTS,
		Nodes: nodes,
	}

	err = pm.SendDataToPeer(op, syncOplogAck, peer)
	if err != nil {
		return err
	}

	return nil
}

func (pm *BaseProtocolManager) HandleSyncOplogAck(
	dataBytes []byte,
	merkle *Merkle,
	syncOplogNewOplogs func(
		syncOplogAck *SyncOplogAck,
		myNewKeys [][]byte,
		theirNewKeys [][]byte,
		peer *PttPeer) error,
	peer *PttPeer) error {

	ptt := pm.Ptt()
	myInfo := ptt.GetMyEntity()
	if myInfo.GetStatus() != types.StatusAlive {
		return nil
	}

	e := pm.Entity()
	if e.GetStatus() != types.StatusAlive {
		return nil
	}

	data := &SyncOplogAck{}
	err := json.Unmarshal(dataBytes, data)
	if err != nil {
		return err
	}

	now, err := types.GetTimestamp()
	if err != nil {
		return err
	}

	myNodes, err := merkle.GetMerkleTreeListByLevel(MerkleTreeLevelNow, data.TS, now)
	if err != nil {
		return err
	}

	myNewKeys, theirNewKeys, err := MergeMerkleNodeKeys(myNodes, data.Nodes)
	if err != nil {
		return err
	}

	return syncOplogNewOplogs(data, myNewKeys, theirNewKeys, peer)
}
