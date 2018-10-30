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

	"github.com/ailabstw/go-pttai/account"
	"github.com/ailabstw/go-pttai/content"
	"github.com/ailabstw/go-pttai/log"
	pkgservice "github.com/ailabstw/go-pttai/service"
)

type PrivateAPI struct {
	b *Backend
}

func NewPrivateAPI(b *Backend) *PrivateAPI {
	return &PrivateAPI{b}
}

func (api *PrivateAPI) SetMyName(name []byte) (*account.UserName, error) {
	return api.b.SetMyName(name)
}

func (api *PrivateAPI) SetMyNodeName(nodeID string, name []byte) (*MyNode, error) {
	return api.b.SetMyNodeName([]byte(nodeID), name)
}

func (api *PrivateAPI) SetMyImage(imgStr string) (*account.UserImg, error) {
	return api.b.SetMyImage(imgStr)
}

func (api *PrivateAPI) Revoke(keyBytes []byte) error {
	return api.b.Revoke(keyBytes)
}

func (api *PrivateAPI) ShowMyKey() (string, error) {
	return api.b.ShowMyKey()
}

func (api *PrivateAPI) ShowMeURL() (*pkgservice.BackendJoinURL, error) {
	return api.b.ShowMeURL()
}

func (api *PrivateAPI) JoinMe(meURL string, myKeyHex string, dummy bool) (*pkgservice.BackendJoinRequest, error) {
	myKeyBytes, err := hex.DecodeString(myKeyHex)
	log.Debug("JoinMe: start", "myKeyHex", myKeyHex, "myKeyBytes", myKeyBytes, "e", err)
	if err != nil {
		log.Error("unable to decode myKey", "myKeyHex", myKeyHex, "e", err)
		return nil, err
	}

	return api.b.JoinMe([]byte(meURL), myKeyBytes)
}

func (api *PrivateAPI) GetMyBoard() (*content.BackendGetBoard, error) {
	return api.b.GetMyBoard()
}

func (api *PrivateAPI) GetMyNodes() ([]*MyNode, error) {
	return api.b.GetMyNodes()
}

func (api *PrivateAPI) GetTotalWeight() uint32 {
	return api.b.GetTotalWeight()
}

func (api *PrivateAPI) GetRawMe() (*MyInfo, error) {
	return api.b.GetRawMe()
}

func (api *PrivateAPI) GetJoinKeyInfos() ([]*pkgservice.KeyInfo, error) {
	idBytes, err := Me.ID.MarshalText()
	if err != nil {
		return nil, err
	}
	return api.b.GetJoinKeyInfos(idBytes)
}

func (api *PrivateAPI) GetOpKeyInfos() ([]*pkgservice.KeyInfo, error) {
	idBytes, err := Me.ID.MarshalText()
	if err != nil {
		return nil, err
	}
	return api.b.GetOpKeyInfos(idBytes)
}

func (api *PrivateAPI) JoinFriend(friendURL string) (*pkgservice.BackendJoinRequest, error) {
	return api.b.JoinFriend([]byte(friendURL))
}

/*
GetFriendRequests get the friend-requests from me to the others.
*/
func (api *PrivateAPI) GetFriendRequests() ([]*pkgservice.BackendJoinRequest, error) {
	return api.b.GetFriendRequests()
}

/*
GetMeRequests get the me-requests from me to the others.
*/
func (api *PrivateAPI) GetMeRequests() ([]*pkgservice.BackendJoinRequest, error) {
	return api.b.GetMeRequests()
}

func (api *PrivateAPI) CountPeers() (int, error) {
	return api.b.CountPeers()
}

func (api *PrivateAPI) GetPeers() ([]*pkgservice.BackendPeer, error) {
	return api.b.GetPeers()
}

func (api *PrivateAPI) GetRaftStatus(id string) (*RaftStatus, error) {
	return api.b.GetRaftStatus([]byte(id))
}

func (api *PrivateAPI) ForceRemoveNode(nodeID string) (bool, error) {
	return api.b.ForceRemoveNode(nodeID)
}

type PublicAPI struct {
	b *Backend
}

func NewPublicAPI(b *Backend) *PublicAPI {
	return &PublicAPI{b}
}

func (api *PublicAPI) Get() (*BackendMyInfo, error) {
	return api.b.Get()
}

func (api *PublicAPI) ShowURL() (*pkgservice.BackendJoinURL, error) {
	return api.b.ShowURL()
}
