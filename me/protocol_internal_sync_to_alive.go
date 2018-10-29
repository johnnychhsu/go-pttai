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
	"github.com/ailabstw/go-pttai/common/types"
	pkgservice "github.com/ailabstw/go-pttai/service"
)

func (pm *ProtocolManager) InternalSyncToAlive(oplog *pkgservice.MasterOplog, weight uint32) error {

	// my-info
	myInfo := pm.Entity().(*MyInfo)
	myInfo.Status = types.StatusAlive
	myInfo.LogID = oplog.ID
	myInfo.UpdateTS = oplog.UpdateTS

	err := myInfo.Save()
	if err != nil {
		return err
	}

	myNode := pm.MyNodes[MyRaftID]
	myNode.Status = types.StatusAlive
	myNode.UpdateTS = oplog.UpdateTS

	_, err = myNode.Save()
	if err != nil {
		return err
	}

	expectedWeight := pm.nodeTypeToWeight(MyNodeType)
	if weight != expectedWeight {
		pm.ProposeRaftAddNode(MyNodeID, expectedWeight)
	}

	return nil
}
