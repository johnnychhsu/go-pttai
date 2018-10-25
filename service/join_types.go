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
	"crypto/ecdsa"

	"github.com/ailabstw/go-pttai/common"
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/p2p/discover"
)

// join-request
type JoinStatus int

const (
	_ JoinStatus = iota
	JoinStatusPending
	JoinStatusRequested
	JoinStatusWaitAccepted
	JoinStatusAccepted
)

type JoinRequest struct {
	CreatorID   *types.PttID      `json:"CID"`
	CreateTS    types.Timestamp   `json:"CT"`
	NodeID      *discover.NodeID  `json:"NID"`
	Hash        *common.Address   `json:"H"`
	Key         *ecdsa.PrivateKey `json:"K"`
	Name        []byte            `json:"N"`
	Status      JoinStatus        `json:"S"`
	Master0Hash []byte            `json:"M"`

	Challenge []byte `json:"C"`
	ID        *types.PttID
}

type JoinRequestEvent struct {
	CreatorID *types.PttID `json:"CID"`
	Hash      []byte       `json:"H"`
	MyID      *types.PttID `json:"MID"`
}

func JoinRequestToJoinRequestEvent(joinRequst *JoinRequest, myID *types.PttID) *JoinRequestEvent {
	return &JoinRequestEvent{
		CreatorID: joinRequst.CreatorID,
		Hash:      joinRequst.Hash[:],
		MyID:      myID,
	}
}

// confirm-join
type JoinType uint8

const (
	JoinTypeInvalid JoinType = iota
	JoinTypeMe
	JoinTypeFriend
	JoinTypeBoard
)

type ConfirmJoin struct {
	Entity     Entity
	JoinEntity *JoinEntity
	KeyInfo    *KeyInfo
	Peer       *PttPeer
	UpdateTS   types.Timestamp
	JoinType   JoinType
}

type JoinEntity struct {
        ID          *types.PttID
        Name        []byte `json:"N"`
        Master0Hash []byte `json:"M"`
}
