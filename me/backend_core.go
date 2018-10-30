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
	"encoding/hex"
	"reflect"

	"github.com/ailabstw/go-pttai/account"
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/content"
	"github.com/ailabstw/go-pttai/crypto"
	"github.com/ailabstw/go-pttai/log"
	"github.com/ailabstw/go-pttai/p2p/discover"
	pkgservice "github.com/ailabstw/go-pttai/service"
	"github.com/syndtr/goleveldb/leveldb"
)

func (b *Backend) SetMyName(name []byte) (*account.UserName, error) {
	return nil, types.ErrNotImplemented
}

func (b *Backend) SetMyNodeName(nodeIDBytes []byte, name []byte) (*MyNode, error) {
	return nil, types.ErrNotImplemented
}

func (b *Backend) SetMyImage(imgStr string) (*account.UserImg, error) {
	return nil, types.ErrNotImplemented
}

func (b *Backend) ShowMeURL() (*pkgservice.BackendJoinURL, error) {
	myInfo := b.SPM().(*ServiceProtocolManager).MyInfo
	pm := myInfo.PM().(*ProtocolManager)

	keyInfo, err := pm.GetJoinKey()
	if err != nil {
		return nil, err
	}

	myUserName := &account.UserName{}
	err = myUserName.Get(myInfo.ID, true)
	if err == leveldb.ErrNotFound {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return pkgservice.MarshalBackendJoinURL(myInfo.ID, MyNodeID, keyInfo, myUserName.Name, pkgservice.PathJoinMe)
}

func (b *Backend) ShowMyKey() (string, error) {
	myKeyBytes := crypto.FromECDSA(MyKey)
	myKeyHex := hex.EncodeToString(myKeyBytes)

	return myKeyHex, nil
}

func (b *Backend) JoinMe(meURL []byte, myKeyBytes []byte) (*pkgservice.BackendJoinRequest, error) {

	joinRequest, err := pkgservice.ParseBackendJoinURL(meURL, pkgservice.PathJoinMe)
	log.Debug("JoinMe: after parse", "joinRequest", joinRequest, "e", err)
	if err != nil {
		return nil, err
	}

	myNodeID := b.Ptt().MyNodeID
	log.Debug("JoinMe: after parse", "joinRequest", joinRequest, "myNodeID", myNodeID, "joinNodeID", joinRequest.NodeID)
	if reflect.DeepEqual(myNodeID, joinRequest.NodeID) {
		return nil, ErrInvalidNode
	}

	myInfo := b.SPM().(*ServiceProtocolManager).MyInfo
	pm := myInfo.PM().(*ProtocolManager)
	err = pm.JoinMe(joinRequest, myKeyBytes)
	if err != nil {
		return nil, err
	}

	backendJoinRequest := pkgservice.JoinRequestToBackendJoinRequest(joinRequest)

	return backendJoinRequest, nil
}

func (b *Backend) ShowURL() (*pkgservice.BackendJoinURL, error) {
	return nil, types.ErrNotImplemented
}

func (b *Backend) JoinFriend(friendURL []byte) (*pkgservice.BackendJoinRequest, error) {
	return nil, types.ErrNotImplemented
}

func (b *Backend) Get() (*BackendMyInfo, error) {
	return MarshalBackendMyInfo(Me), nil
}

func (b *Backend) GetRawMe() (*MyInfo, error) {
	return Me, nil
}

func (b *Backend) Revoke(keyBytes []byte) error {
	return nil
}

func (b *Backend) GetMyNodes() ([]*MyNode, error) {
	pm := b.SPM().(*ServiceProtocolManager).MyInfo.PM().(*ProtocolManager)

	pm.RLockMyNodes()
	defer pm.RUnlockMyNodes()

	myNodeList := make([]*MyNode, len(pm.MyNodes))

	i := 0
	for _, node := range pm.MyNodes {
		myNodeList[i] = node

		i++
	}

	return myNodeList, nil
}

func (b *Backend) GetFriendRequests() ([]*pkgservice.BackendJoinRequest, error) {
	return nil, types.ErrNotImplemented
}

func (b *Backend) GetMeRequests() ([]*pkgservice.BackendJoinRequest, error) {
	pm := b.SPM().(*ServiceProtocolManager).MyInfo.PM().(*ProtocolManager)

	joinMeRequests, lockJoinMeRequest := pm.GetJoinMeRequests()

	lockJoinMeRequest.RLock()
	defer lockJoinMeRequest.RUnlock()

	lenRequests := len(joinMeRequests)
	results := make([]*pkgservice.BackendJoinRequest, lenRequests)

	i := 0
	for _, joinRequest := range joinMeRequests {
		results[i] = pkgservice.JoinRequestToBackendJoinRequest(joinRequest)

		i++
	}

	return results, nil
}

func (b *Backend) CountPeers() (int, error) {
	pm := b.SPM().(*ServiceProtocolManager).MyInfo.PM().(*ProtocolManager)

	return pm.CountPeers()
}

func (b Backend) GetPeers() ([]*pkgservice.BackendPeer, error) {
	pm := b.SPM().(*ServiceProtocolManager).MyInfo.PM().(*ProtocolManager)

	peerList, err := pm.GetPeers()
	if err != nil {
		return nil, err
	}

	backendPeerList := make([]*pkgservice.BackendPeer, len(peerList))

	for i, peer := range peerList {
		backendPeerList[i] = pkgservice.PeerToBackendPeer(peer)
	}

	return backendPeerList, nil
}

func (b *Backend) GetMyBoard() (*content.BackendGetBoard, error) {
	return nil, types.ErrNotImplemented
}

func (b *Backend) GetRaftStatus(idBytes []byte) (*RaftStatus, error) {
	var myInfo *MyInfo
	if len(idBytes) == 0 {
		myInfo = b.SPM().(*ServiceProtocolManager).MyInfo
	} else {
		myID := &types.PttID{}
		err := myID.UnmarshalText(idBytes)
		if err != nil {
			return nil, err
		}

		myInfo = b.SPM().Entity(myID).(*MyInfo)
	}

	if myInfo == nil {
		return nil, ErrInvalidMe
	}

	pm := myInfo.PM().(*ProtocolManager)
	return pm.GetRaftStatus()
}

func (b *Backend) GetTotalWeight() uint32 {
	myInfo := b.SPM().(*ServiceProtocolManager).MyInfo
	pm := myInfo.PM().(*ProtocolManager)

	return pm.totalWeight
}

func (b *Backend) ForceRemoveNode(nodeIDStr string) (bool, error) {
	myInfo := b.SPM().(*ServiceProtocolManager).MyInfo

	nodeID, err := discover.HexID(nodeIDStr)
	if err != nil {
		return false, err
	}

	pm := myInfo.PM().(*ProtocolManager)
	err = pm.ForceProposeRaftRemoveNode(&nodeID)
	return false, err
}
