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

type MeOplog struct {
	*BaseOplog `json:"O"`
}

func NewMeOplog(objID *types.PttID, ts types.Timestamp, doerID *types.PttID, op OpType, data interface{}, myID *types.PttID) (*MeOplog, error) {

	oplog, err := NewOplog(objID, ts, doerID, op, data, dbOplog, myID, DBMeOplogPrefix, DBMeIdxOplogPrefix, DBMeMerkleOplogPrefix, DBMeLockMap)
	if err != nil {
		return nil, err
	}

	return &MeOplog{
		BaseOplog: oplog,
	}, nil
}

func SetMeOplogDB(meOplog *MeOplog, myID *types.PttID) {
	meOplog.SetDB(dbOplog, myID, DBMeOplogPrefix, DBMeIdxOplogPrefix, DBMeMerkleOplogPrefix, DBMeLockMap)
}
