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

package me

import (
	"encoding/json"

	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/log"
	pkgservice "github.com/ailabstw/go-pttai/service"
)

type InitMeInfoAck struct {
	Status types.Status `json:"S"`
}

func (pm *ProtocolManager) InitMeInfoAck(data *InitMeInfo, peer *pkgservice.PttPeer) error {

	ts, err := types.GetTimestamp()
	if err != nil {
		return err
	}

	myInfo := pm.Entity().(*MyInfo)

	if myInfo.Status == types.StatusInit {
		myInfo.Status = types.StatusInternalPending
		myInfo.UpdateTS = ts
		err = myInfo.Save()
		if err != nil {
			return err
		}

		myNode := pm.MyNodes[MyRaftID]
		myNode.Status = myInfo.Status
		myNode.UpdateTS = ts
		_, err = myNode.Save()
		if err != nil {
			return err
		}
	}
	log.Debug("InitMeInfoAck: start", "myID", myInfo.ID, "Status", myInfo.Status)

	return pm.SendDataToPeer(InitMeInfoAckMsg, &InitMeInfoAck{Status: myInfo.Status}, peer)

}

func (pm *ProtocolManager) HandleInitMeInfoAck(dataBytes []byte, peer *pkgservice.PttPeer) error {
	data := &InitMeInfoAck{}
	err := json.Unmarshal(dataBytes, data)
	if err != nil {
		return err
	}

	if data.Status == types.StatusInternalPending {
		return pm.InitMeInfoSync(peer)
	}

	nodeID := peer.GetID()
	raftID, err := nodeID.ToRaftID()
	if err != nil {
		return err
	}

	myNode := pm.MyNodes[raftID]
	if myNode == nil {
		return ErrInvalidNode
	}

	if myNode.Status == data.Status {
		return nil
	}

	log.Debug("HandleInitMeInfoAck: start", "peerID", nodeID, "status", data.Status)

	ts, err := types.GetTimestamp()
	if err != nil {
		return err
	}

	myNode.Status = data.Status
	myNode.UpdateTS = ts
	_, err = myNode.Save()
	if err != nil {
		return err
	}

	return nil
}
