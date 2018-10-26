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

// ISC License
//
// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2016-2017 The Lightning Network Developers
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package service

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"

	"github.com/ailabstw/go-pttai/common"
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/crypto"
	"github.com/ailabstw/go-pttai/log"
	"github.com/ailabstw/go-pttai/pttdb"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/pbkdf2"
)

// KeyInfo

type KeyInfo struct {
	ID       *types.PttID
	Key      *ecdsa.PrivateKey `json:"-"`
	KeyBytes []byte            `json:"K"`
	Hash     *common.Address   `json:"H"`
	UpdateTS types.Timestamp   `json:"UT"`
	EntityID *types.PttID      `json:"EID"`
	DoerID   *types.PttID      `json:"DID"`
	Status   types.Status      `json:"S"`
	LogID    *types.PttID      `json:"pl"`
	Extra    *KeyExtraInfo     `json:"e,omitempty"`
}

func NewOpKeyInfo(entityID *types.PttID, doerID *types.PttID, masterKey *ecdsa.PrivateKey) (*KeyInfo, error) {
	key, extra, err := deriveOpKey(masterKey)
	if err != nil {
		return nil, err
	}
	hash := crypto.PubkeyToAddress(key.PublicKey)

	updateTS, err := types.GetTimestamp()
	if err != nil {
		return nil, err
	}

	id := &types.PttID{}
	copy(id[:common.AddressLength], hash[:])

	keyBytes := crypto.FromECDSA(key)

	return &KeyInfo{
		ID:       id,
		Key:      key,
		KeyBytes: keyBytes,
		Hash:     &hash,
		UpdateTS: updateTS,
		EntityID: entityID,
		Extra:    extra,
	}, nil
}

func NewJoinKeyInfo(entityID *types.PttID, doerID *types.PttID) (*KeyInfo, error) {
	key, err := deriveJoinKey()
	if err != nil {
		return nil, err
	}

	return newKeyInfo(key, nil, entityID, doerID)
}

/*
NewSignKeyInfo is used for signing for now.
*/
func NewSignKeyInfo(entityID *types.PttID, doerID *types.PttID, masterKey *ecdsa.PrivateKey) (*KeyInfo, error) {
	key, extra, err := deriveSignKey(masterKey)
	if err != nil {
		return nil, err
	}

	return newKeyInfo(key, extra, entityID, doerID)
}

func newKeyInfo(key *ecdsa.PrivateKey, extra *KeyExtraInfo, entityID *types.PttID, doerID *types.PttID) (*KeyInfo, error) {
	hash := crypto.PubkeyToAddress(key.PublicKey)

	updateTS, err := types.GetTimestamp()
	if err != nil {
		return nil, err
	}

	id := &types.PttID{}
	copy(id[:common.AddressLength], hash[:])

	keyBytes := crypto.FromECDSA(key)

	return &KeyInfo{
		ID:       id,
		Key:      key,
		KeyBytes: keyBytes,
		Hash:     &hash,
		UpdateTS: updateTS,
		EntityID: entityID,
		DoerID:   doerID,
		Extra:    extra,
	}, nil
}

func deriveJoinKey() (*ecdsa.PrivateKey, error) {
	key, err := crypto.GenerateKey()
	return key, err
}

func deriveOpKey(masterKey *ecdsa.PrivateKey) (*ecdsa.PrivateKey, *KeyExtraInfo, error) {
	return deriveKeyBIP32(masterKey)
}

func deriveSignKey(masterKey *ecdsa.PrivateKey) (*ecdsa.PrivateKey, *KeyExtraInfo, error) {
	return deriveKeyBIP32(masterKey)
}

func deriveKeyBIP32(masterKey *ecdsa.PrivateKey) (*ecdsa.PrivateKey, *KeyExtraInfo, error) {
	salt, err := types.NewSalt()
	if err != nil {
		return nil, nil, err
	}
	saltBytes := salt[:]
	masterKeyBytes := crypto.FromECDSA(masterKey)
	keyBytes := pbkdf2.Key(masterKeyBytes, saltBytes, 1, 32, sha256.New)

	key, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, nil, err
	}

	keyData := &KeyBIP32{
		Salt: saltBytes,
	}

	extra := &KeyExtraInfo{
		KeyType: KeyTypeBIP32,
		Data:    keyData,
	}

	return key, extra, err
}

func (k *KeyInfo) Save(db *pttdb.LDBBatch) error {
	var err error
	if k.Key == nil && k.KeyBytes != nil {
		k.Key, err = crypto.ToECDSA(k.KeyBytes)
		if err != nil {
			return err
		}
	}
	key, err := k.MarshalKey()
	if err != nil {
		return err
	}
	marshaled, err := k.Marshal()
	if err != nil {
		return err
	}

	idxKey, err := k.IdxKey()
	if err != nil {
		return err
	}

	idxKey2, err := k.IdxKey2()
	if err != nil {
		return err
	}

	idx := &pttdb.Index{Keys: [][]byte{key}, UpdateTS: k.UpdateTS}

	kvs := []*pttdb.KeyVal{
		&pttdb.KeyVal{K: key, V: marshaled},
		&pttdb.KeyVal{K: idxKey2, V: key},
	}

	_, err = db.TryPutAll(idxKey, idx, kvs, true, false)
	if err != nil {
		return err
	}

	return nil
}

func (k *KeyInfo) LoadNewest(db *pttdb.LDBBatch) error {
	iter, err := db.DB().NewPrevIteratorWithPrefix(nil, k.DBPrefix())
	if err != nil {
		return err
	}
	defer iter.Release()

	log.Debug("LoadNewest: to loop for", "DBPrefix", k.DBPrefix(), "iter", iter)

	var val []byte = nil
	for iter.Prev() {
		val = iter.Value()
		log.Debug("LoadNewest: in-loop", "val", val)
		break

	}
	if val == nil {
		return leveldb.ErrNotFound
	}

	err = k.Unmarshal(val)
	if err != nil {
		return err
	}

	k.Key, err = crypto.ToECDSA(k.KeyBytes)
	if err != nil {
		return err
	}

	return nil
}

func (k *KeyInfo) IdxPrefix() []byte {
	return append(DBOpKeyIdxPrefix, k.EntityID[:]...)
}

func (k *KeyInfo) IdxKey() ([]byte, error) {
	return common.Concat([][]byte{DBOpKeyIdxPrefix, k.EntityID[:], k.ID[:]})
}

func (k *KeyInfo) IdxKey2() ([]byte, error) {
	return common.Concat([][]byte{DBOpKeyIdx2Prefix, k.EntityID[:], k.Hash[:]})
}

func (k *KeyInfo) DBPrefix() []byte {
	return append(DBOpKeyPrefix, k.EntityID[:]...)
}

func (k *KeyInfo) MarshalKey() ([]byte, error) {
	marshalTimestamp, err := k.UpdateTS.Marshal()
	if err != nil {
		return nil, err
	}
	return common.Concat([][]byte{DBOpKeyPrefix, k.EntityID[:], marshalTimestamp, k.ID[:]})
}

func (k *KeyInfo) KeyToIdxKey(key []byte) ([]byte, error) {

	lenKey := len(key)
	if lenKey != pttdb.SizeDBKeyPrefix+types.SizePttID+types.SizeTimestamp+types.SizePttID {
		return nil, ErrInvalidKey
	}

	idxKey := make([]byte, pttdb.SizeDBKeyPrefix+types.SizePttID+types.SizePttID)

	// prefix
	idxOffset := 0
	nextIdxOffset := pttdb.SizeDBKeyPrefix
	copy(idxKey[:nextIdxOffset], DBOpKeyIdxPrefix)

	// entity-id
	idxOffset = nextIdxOffset
	nextIdxOffset += types.SizePttID

	keyOffset := pttdb.SizeDBKeyPrefix
	nextKeyOffset := keyOffset + types.SizePttID
	copy(idxKey[idxOffset:nextIdxOffset], key[keyOffset:nextKeyOffset])

	// id
	idxOffset = nextIdxOffset
	nextIdxOffset += types.SizePttID

	keyOffset = lenKey - types.SizePttID
	nextKeyOffset = lenKey
	copy(idxKey[idxOffset:nextIdxOffset], key[keyOffset:nextKeyOffset])

	return idxKey, nil
}

func (k *KeyInfo) DeleteKey(key []byte, db *pttdb.LDBBatch) error {
	idxKey, err := k.KeyToIdxKey(key)
	if err != nil {
		return err
	}

	log.Debug("DeleteKey", "idxKey", idxKey)

	err = db.DeleteAll(idxKey)

	if err != nil {
		return err
	}

	return nil
}

func (k *KeyInfo) Marshal() ([]byte, error) {
	return json.Marshal(k)
}

func (k *KeyInfo) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, k)
	if err != nil {
		return err
	}

	k.Key, err = crypto.ToECDSA(k.KeyBytes)
	if err != nil {
		return err
	}

	return nil
}
