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
	"bytes"
	"os"

	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/content"
	"github.com/ailabstw/go-pttai/log"
	pkgservice "github.com/ailabstw/go-pttai/service"
)

func (spm *ServiceProtocolManager) CreateMe(contentBackend *content.Backend) error {
	ptt := spm.myPtt

	// new my info
	myInfo, err := NewMyInfo(MyID, MyKey, ptt, spm.Service())
	if err != nil {
		return err
	}

	// save
	err = myInfo.Save()
	if err != nil {
		return err
	}

	// add to entities
	err = spm.RegisterEntity(myInfo.ID, myInfo)
	if err != nil {
		return err
	}

	// Me
	Me = myInfo

	// SPM MyInfo
	spm.MyInfo = Me

	return nil
}

func (pm *ProtocolManager) CreateFullMe(oplog *pkgservice.MasterOplog) error {
	log.Debug("CreateFullMe: start")
	myInfo := pm.Entity().(*MyInfo)
	ptt := pm.myPtt

	// create-me-oplog

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	myHostname := bytes.TrimSuffix([]byte(hostname), []byte(".local"))

	opData := &pkgservice.MeOpCreateMe{
		NodeID:   MyNodeID,
		NodeType: MyNodeType,
		NodeName: myHostname,
	}

	meOplog, err := ptt.CreateMeOplog(MyID, oplog.UpdateTS, pkgservice.MeOpTypeCreateMe, opData)
	if err != nil {
		return err
	}
	log.Debug("CreateFullMe: after CreateMeOplog", "masterLogID", meOplog.MasterLogID)

	// my-info
	myInfo.Status = types.StatusAlive
	myInfo.CreateTS = meOplog.UpdateTS
	myInfo.UpdateTS = meOplog.UpdateTS
	myInfo.LogID = meOplog.ID

	err = myInfo.Save()
	if err != nil {
		return err
	}

	// my-node
	myNode, ok := pm.MyNodes[MyRaftID]
	if !ok {
		return ErrInvalidNode
	}

	myNode.NodeName = myHostname
	myNode.Save()

	// meOplog save
	meOplog.Save(false)

	return nil
}
