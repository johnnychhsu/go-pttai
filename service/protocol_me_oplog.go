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

func (p *BasePtt) CreateMeOplog(objID *types.PttID, ts types.Timestamp, op OpType, data interface{}) (*MeOplog, error) {
	myEntity := p.myEntity

	myID := myEntity.GetID()

	key := myEntity.SignKey()
	nodeSignID := myEntity.GetNodeSignID()

	oplog, err := NewMeOplog(objID, ts, nodeSignID, op, data, myID)
	if err != nil {
		return nil, err
	}
	err = oplog.Sign(key)
	if err != nil {
		return nil, err
	}

	err = p.MasterSignMeOplog(oplog)
	if err != nil {
		return nil, err
	}

	return oplog, nil
}

func (p *BasePtt) MasterSignMeOplog(oplog *MeOplog) error {
	if oplog.MasterLogID != nil {
		return nil
	}

	myEntity := p.myEntity
	key := myEntity.SignKey()
	nodeSignID := myEntity.GetNodeSignID()

	err := oplog.MasterSign(nodeSignID, key)
	if err != nil {
		return err
	}

	masterLogID, weight, isValid := myEntity.MyPM().IsValidOplog(oplog.MasterSigns)
	if !isValid {
		return nil
	}

	err = oplog.SetMasterLogID(masterLogID, weight)
	if err != nil {
		return err
	}

	return nil
}

func (p *BasePtt) GetMeOplogsFromKeys(keys [][]byte) ([]*MeOplog, error) {
	pm := p.myEntity.MyPM()
	logs, err := pm.GetOplogsFromKeys(pm.SetMeDB, keys)
	if err != nil {
		return nil, err
	}

	opKeyLogs := OplogsToMeOplogs(logs)

	return opKeyLogs, nil
}

func (p *BasePtt) IntegrateMeOplog(log *MeOplog, isLocked bool) (bool, error) {
	pm := p.myEntity.MyPM()
	return pm.IntegrateOplog(log.Oplog, isLocked)
}

func (p *BasePtt) GetPendingMeOplogs() ([]*MeOplog, []*MeOplog, error) {
	pm := p.myEntity.MyPM()
	logs, failedLogs, err := pm.GetPendingOplogs(pm.SetMeDB)
	if err != nil {
		return nil, nil, err
	}

	opKeyLogs := OplogsToMeOplogs(logs)

	failedMeLogs := OplogsToMeOplogs(failedLogs)

	return opKeyLogs, failedMeLogs, nil
}

func (p *BasePtt) BroadcastMeOplog(log *MeOplog) error {
	pm := p.myEntity.MyPM()
	return pm.BroadcastOplog(log.Oplog, AddMeOplogMsg, AddPendingMeOplogMsg)
}

func (p *BasePtt) BroadcastMeOplogs(opKeyLogs []*MeOplog) error {
	pm := p.myEntity.MyPM()
	logs := MeOplogsToOplogs(opKeyLogs)
	return pm.BroadcastOplogs(logs, AddMeOplogsMsg, AddPendingMeOplogsMsg)
}

func (p *BasePtt) SetMeOplogIsSync(log *MeOplog, isBroadcast bool) (bool, error) {
	pm := p.myEntity.MyPM()
	isNewSign, err := pm.SetOplogIsSync(log.Oplog)
	if err != nil {
		return false, err
	}
	if isNewSign && isBroadcast {
		p.BroadcastMeOplog(log)
	}

	return isNewSign, nil
}

func (p *BasePtt) RemoveNonSyncMeOplog(logID *types.PttID, isRetainValid bool, isLocked bool) (*MeOplog, error) {

	pm := p.myEntity.MyPM()
	log, err := pm.RemoveNonSyncOplog(pm.SetMeDB, logID, isRetainValid, isLocked)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, nil
	}

	return &MeOplog{Oplog: log}, nil
}
