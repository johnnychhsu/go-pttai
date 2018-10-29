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
	"context"
	"encoding/json"

	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/log"
	pb "github.com/ailabstw/go-pttai/raft/raftpb"
	pkgservice "github.com/ailabstw/go-pttai/service"
)

type SendRaftMsgs struct {
	Msgs []pb.Message `json:"M"`
}

func (pm *ProtocolManager) SendRaftMsgs(msgs []pb.Message) error {
	log.Debug("SendRaftMsgs: start", "msgs", msgs)
	msgByPeers := make(map[uint64][]pb.Message)
	var origMsgByPeers []pb.Message
	for _, msg := range msgs {
		if msg.To == MyRaftID {
			log.Debug("SendRaftMsgs: To is myRaftID", "msg", msg, "MyRaftID", MyRaftID)
			continue
		}
		origMsgByPeers = msgByPeers[msg.To]
		msgByPeers[msg.To] = append(origMsgByPeers, msg)
	}

	pm.LockMyNodes.RLock()
	defer pm.LockMyNodes.RUnlock()

	peers := pm.Peers().Peers()
	var data *SendRaftMsgs
	for raftID, eachMsgs := range msgByPeers {
		myNode := pm.MyNodes[raftID]
		if myNode == nil {
			log.Debug("SendRaftMsgs: unable to send peer not myNode", "raftID", raftID)
			continue
		}

		if myNode.Status == types.StatusInit || myNode.Status == types.StatusInternalPending {
			log.Debug("SendRaftMsgs: myNode status invalid", "raftID", raftID, "status", myNode.Status)
			continue
		}

		peer := peers[*myNode.NodeID]
		if peer == nil {
			log.Debug("SendRaftMsgs: unable to send peer", "nodeID", myNode.NodeID)
			continue
		}

		data = &SendRaftMsgs{
			Msgs: eachMsgs,
		}

		pkgservice.PMSendDataToPeer(pm, SendRaftMsgsMsg, data, peer)
	}

	return nil
}

func (pm *ProtocolManager) HandleSendRaftMsgs(dataBytes []byte, peer *pkgservice.PttPeer) error {
	myInfo := pm.Entity().(*MyInfo)
	// defensive-programming
	if myInfo.Status == types.StatusInit || myInfo.Status == types.StatusInternalPending {
		return nil
	}

	if pm.raftNode == nil {
		return nil
	}

	// unmarshal
	data := &SendRaftMsgs{}
	err := json.Unmarshal(dataBytes, data)
	log.Debug("HandleSendRaftMsgs: start", "myID", myInfo.ID, "peer", peer.GetID(), "msgs", data.Msgs)
	if err != nil {
		return err
	}

	// step
	for _, msg := range data.Msgs {
		pm.raftNode.Step(context.TODO(), msg)
	}

	return nil
}
