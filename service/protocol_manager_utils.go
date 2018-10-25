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
	"encoding/json"
	"time"

	"github.com/ailabstw/go-pttai/common"
	"github.com/ailabstw/go-pttai/log"
	"github.com/ailabstw/go-pttai/p2p"
)

func StartPM(pm ProtocolManager) error {
	entityName := pm.Entity().Name()
	log.Info("StartPM: start", "entity", entityName)

	go func() {
		pm.SyncWG().Add(1)
		defer func() {
			log.Debug("StartPM: finish sync. to syncWG.Done", "entity", entityName)
			pm.SyncWG().Done()
		}()

		PMSync(pm)
	}()

	err := pm.Start()
	if err != nil {
		return err
	}

	return nil
}

func StopPM(pm ProtocolManager) error {
	entityName := pm.Entity().Name()
	log.Info("Stop PM: to stop", "entity", entityName)

	err := pm.Stop()
	if err != nil {
		return err
	}

	log.Info(entityName + " protocol stopped")

	return nil
}

func PMSync(pm ProtocolManager) error {
	var err error

	pm.SetForceSyncCycle()
	forceSyncTicker := time.NewTicker(pm.ForceSyncCycle())

looping:
	for {
		select {
		case p, ok := <-pm.NewPeerCh():
			if !ok {
				break looping
			}

			err = pm.Sync(p)
			if err != nil {
				log.Error("unable to Sync after newPeer", "e", err)
			}
		case <-forceSyncTicker.C:
			forceSyncTicker.Stop()
			pm.SetForceSyncCycle()
			forceSyncTicker = time.NewTicker(pm.ForceSyncCycle())

			err = pm.Sync(nil)
			if err != nil {
				log.Error("unable to Sync after forceSync", "e", err)
			}
		case <-pm.QuitSync():
			return p2p.DiscQuitting
		}
	}

	forceSyncTicker.Stop()

	return nil
}

func PMHandleMessageWrapper(pm ProtocolManager, hash *common.Address, encData []byte, peer *PttPeer) error {
	opKeyInfo, err := pm.GetOpKeyInfo(hash)

	if err != nil {
		return err
	}

	op, dataBytes, err := pm.Ptt().DecryptData(encData, opKeyInfo.Key)
	if err != nil {
		return err
	}

	fitPeerType := pm.IsFitPeer(peer)

	if fitPeerType < PeerTypeMember {
		return ErrInvalidEntity
	}

	return pm.HandleMessage(op, dataBytes, peer)
}

/*
Send Data to Peers using op-key
*/
func PMSendDataToPeers(pm ProtocolManager, op OpType, data interface{}, peerList []*PttPeer) error {

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	opKeyInfo, err := pm.GetOpKey()
	if err != nil {
		return err
	}

	ptt := pm.Ptt()
	encData, err := ptt.EncryptData(op, dataBytes, opKeyInfo.Key)
	if err != nil {
		return err
	}

	pttData, err := ptt.MarshalData(CodeTypeOp, opKeyInfo.Hash, encData)
	if err != nil {
		return err
	}

	okCount := 0
	for _, peer := range peerList {
		pttData.Node = peer.GetID()[:]
		err := peer.SendData(pttData)
		if err == nil {
			okCount++
		} else {
			log.Warn("PMSendDataToPeers: unable to SendData", "peer", peer, "e", err)
			pm.UnregisterPeer(peer)
		}
	}
	if okCount == 0 {
		return ErrNotSent
	}

	return nil
}

/*
Send Data to Peers using op-key
*/
func PMSendDataToPeer(pm ProtocolManager, op OpType, data interface{}, peer *PttPeer) error {

	// /log.Debug("PMSendDataToPeer: to marshal")
	dataBytes, err := json.Marshal(data)
	//log.Debug("PMSendDataToPeer: after marshal", "e", err)
	if err != nil {
		return err
	}

	opKeyInfo, err := pm.GetOpKey()
	//log.Debug("PMSendDataToPeer: after get ok-key", "e", err)
	if err != nil {
		return err
	}

	ptt := pm.Ptt()
	encData, err := ptt.EncryptData(op, dataBytes, opKeyInfo.Key)
	//log.Debug("PMSendDataToPeer: after enc", "e", err)
	if err != nil {
		return err
	}

	pttData, err := ptt.MarshalData(CodeTypeOp, opKeyInfo.Hash, encData)
	//log.Debug("PMSendDataToPeer: after marshal", "e", err)
	if err != nil {
		return err
	}
	//id := pm.Entity().GetID()
	//name := pm.Entity().Name()
	//log.Debug("PMSendDataToPeer", "pm", pm, "op", op, "id", id, "name", name)

	pttData.Node = peer.GetID()[:]
	err = peer.SendData(pttData)
	//log.Debug("PMSendDataToPeer: after send-data", "e", err)
	if err != nil {
		return ErrNotSent
	}

	return nil
}

func PMPeerListWithType(pm ProtocolManager, peerList []*PttPeer) ([]*PttPeer, []*PttPeer, []*PttPeer) {
	lenPeerList := len(peerList)
	mePeers := make([]*PttPeer, 0, lenPeerList)
	masterPeers := make([]*PttPeer, 0, lenPeerList)
	memberPeers := make([]*PttPeer, 0, lenPeerList)

	for _, peer := range peerList {
		if peer.PeerType == PeerTypeMe {
			mePeers = append(mePeers, peer)
		} else if peer.UserID != nil && pm.IsMaster(peer.UserID) {
			masterPeers = append(masterPeers, peer)
		} else {
			memberPeers = append(memberPeers, peer)
		}
	}

	return mePeers, masterPeers, memberPeers
}
