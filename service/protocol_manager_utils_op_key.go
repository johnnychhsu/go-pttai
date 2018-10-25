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
	"reflect"
	"sort"

	"github.com/ailabstw/go-pttai/common"
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/log"
)

func (pm *BaseProtocolManager) GetOpKeyInfo(hash *common.Address) (*KeyInfo, error) {
	pm.lockOpKeyInfo.RLock()
	defer pm.lockOpKeyInfo.RUnlock()

	var keyInfo *KeyInfo = nil
	for _, eachKeyInfo := range pm.opKeyInfos {
		if reflect.DeepEqual(hash, eachKeyInfo.Hash) {
			keyInfo = eachKeyInfo
			break
		}
	}

	if keyInfo == nil {
		return nil, ErrInvalidKeyInfo
	}

	return keyInfo, nil
}

func (pm *BaseProtocolManager) GetOpKey() (*KeyInfo, error) {
	pm.lockOpKeyInfo.RLock()
	defer pm.lockOpKeyInfo.RUnlock()

	lenKeyInfo := len(pm.opKeyInfos)

	if lenKeyInfo == 0 {
		return nil, ErrInvalidKeyInfo
	}

	return pm.opKeyInfos[lenKeyInfo-1], nil
}

func (pm *BaseProtocolManager) DBOpKeyLock() *types.LockMap {
	return pm.dbOpKeyLock
}

func (pm *BaseProtocolManager) SaveOpKeyInfo(opKeyInfo *KeyInfo) error {
	return opKeyInfo.Save(pm.db)
}

func (pm *BaseProtocolManager) OpKeyInfos() []*KeyInfo {
	return pm.opKeyInfos
}

func (pm *BaseProtocolManager) RegisterOpKeyInfo(keyInfo *KeyInfo) error {
	pm.lockOpKeyInfo.Lock()
	defer pm.lockOpKeyInfo.Unlock()

	opKeyInfos := pm.opKeyInfos
	idx := sort.Search(len(opKeyInfos), func(i int) bool {
		return keyInfo.UpdateTS.IsLess(opKeyInfos[i].UpdateTS)
	})
	pm.opKeyInfos = append(append(opKeyInfos[:idx], keyInfo), opKeyInfos[idx:]...)

	return nil
}

func (pm *BaseProtocolManager) RenewOpKeySeconds() uint64 {
	return pm.renewOpKeySeconds
}

func (pm *BaseProtocolManager) ExpireOpKeySeconds() uint64 {
	return pm.expireOpKeySeconds
}

func (pm *BaseProtocolManager) loadOpKeyInfos(expireOpKeySeconds uint64) ([]*KeyInfo, error) {
	e := pm.Entity()
	entityID := e.GetID()

	dbPrefix, err := DBPrefix(DBOpKeyPrefix, entityID)
	if err != nil {
		return nil, err
	}

	iter, err := pm.db.DB().NewIteratorWithPrefix(nil, dbPrefix)
	if err != nil {
		return nil, err
	}
	defer iter.Release()

	now, err := types.GetTimestamp()
	if err != nil {
		return nil, err
	}
	expireTS := now
	expireTS.Ts -= expireOpKeySeconds

	opKeyInfos := make([]*KeyInfo, 0)
	toRemoveOpKeys := make([][]byte, 0)
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()

		keyInfo := &KeyInfo{}
		err := keyInfo.Unmarshal(val)
		if err != nil {
			log.Warn("loadOpKeyInfo: unable to unmarshal", "key", key, "v", val)
			toRemoveOpKeys = append(toRemoveOpKeys, common.CloneBytes(key))
			continue
		}

		if keyInfo.UpdateTS.IsLess(expireTS) {
			log.Warn("loadOpKeyInfo: expire", "key", key, "expireTS", expireTS, "UpdateTS", keyInfo.UpdateTS)
			toRemoveOpKeys = append(toRemoveOpKeys, common.CloneBytes(key))
			continue
		}

		opKeyInfos = append(opKeyInfos, keyInfo)
	}

	for _, key := range toRemoveOpKeys {
		keyInfo := &KeyInfo{}
		err := keyInfo.DeleteKey(key, pm.db)
		if err != nil {
			log.Error("loadOpKeyInfos: unable to delete key", "name", e.Name(), "dbPrefix", dbPrefix, "e", err)
		}
	}

	return opKeyInfos, nil
}
