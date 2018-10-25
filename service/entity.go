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
)

type Entity interface {
	// variables
	GetVersion() types.Version
	GetID() *types.PttID
	GetCreateTS() types.Timestamp

	GetUpdateTS() types.Timestamp
	SetUpdateTS(ts types.Timestamp)

	GetStatus() types.Status
	SetStatus(status types.Status)

	GetLogID() *types.PttID
	SetLogID(logID *types.PttID)

	GetOwnerID() *types.PttID
	SetOwnerID(id *types.PttID)

	// private-variables

	PM() ProtocolManager
	Ptt() Ptt
	Service() Service

	Name() string
	SetName(name string)

	// start / stop
	Start() error
	Stop() error
}

type BaseEntity struct {
	V         types.Version
	ID        *types.PttID
	CreateTS  types.Timestamp `json:"CT"`
	CreatorID *types.PttID    `json:"CID"`
	UpdateTS  types.Timestamp `json:"UT"`
	UpdaterID *types.PttID    `json:"UID"`

	Status types.Status `json:"S"`

	LogID *types.PttID `json:"l"`

	OwnerID *types.PttID `json:"o,omitempty"`

	pm      ProtocolManager
	ptt     Ptt
	service Service

	name string
}

func NewBaseEntity(id *types.PttID, creatorID *types.PttID, pm ProtocolManager, name string, ptt Ptt, service Service) (*BaseEntity, error) {

	ts, err := types.GetTimestamp()
	if err != nil {
		return nil, err
	}

	b := &BaseEntity{
		V:         types.CurrentVersion,
		ID:        id,
		CreateTS:  ts,
		CreatorID: creatorID,
		UpdateTS:  ts,
		UpdaterID: creatorID,

		pm:      pm,
		name:    name,
		ptt:     ptt,
		service: service,
	}

	return b, nil
}

func (b *BaseEntity) GetVersion() types.Version {
	return b.V
}

func (b *BaseEntity) GetID() *types.PttID {
	return b.ID
}

func (b *BaseEntity) GetCreateTS() types.Timestamp {
	return b.CreateTS
}

func (b *BaseEntity) GetUpdateTS() types.Timestamp {
	return b.UpdateTS
}

func (b *BaseEntity) SetUpdateTS(ts types.Timestamp) {
	b.UpdateTS = ts
}

func (b *BaseEntity) GetStatus() types.Status {
	return b.Status
}

func (b *BaseEntity) SetStatus(status types.Status) {
	b.Status = status
}

func (b *BaseEntity) GetLogID() *types.PttID {
	return b.LogID
}

func (b *BaseEntity) SetLogID(logID *types.PttID) {
	b.LogID = logID
}

func (b *BaseEntity) GetOwnerID() *types.PttID {
	return b.OwnerID
}

func (b *BaseEntity) SetOwnerID(id *types.PttID) {
	b.OwnerID = id
}

func (b *BaseEntity) PM() ProtocolManager {
	return b.pm
}

func (b *BaseEntity) Ptt() Ptt {
	return b.ptt
}

func (b *BaseEntity) Service() Service {
	return b.service
}

func (b *BaseEntity) Name() string {
	return b.name
}

func (b *BaseEntity) SetName(name string) {
	b.name = name
}

func (b *BaseEntity) Start() error {
	return StartPM(b.pm)
}

func (b *BaseEntity) Stop() error {
	return StopPM(b.pm)
}
