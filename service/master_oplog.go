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

import "github.com/ailabstw/go-pttai/common/types"

type MasterOplog struct {
	*Oplog `json:"O"`
}

func NewMasterOplog(id *types.PttID, ts types.Timestamp, doerID *types.PttID, op OpType, data interface{}) (*MasterOplog, error) {

	log, err := NewOplog(id, ts, doerID, op, data, dbOplog, id, DBMasterOplogPrefix, DBMasterIdxOplogPrefix, DBMasterMerkleOplogPrefix, DBMasterLockMap)
	if err != nil {
		return nil, err
	}
	return &MasterOplog{
		Oplog: log,
	}, nil
}

func OplogsToMasterOplogs(logs []*Oplog) []*MasterOplog {
	typedLogs := make([]*MasterOplog, len(logs))
	for i, log := range logs {
		typedLogs[i] = &MasterOplog{Oplog: log}
	}
	return typedLogs
}

func MasterOplogsToOplogs(typedLogs []*MasterOplog) []*Oplog {
	logs := make([]*Oplog, len(typedLogs))
	for i, log := range typedLogs {
		logs[i] = log.Oplog
	}
	return logs
}
