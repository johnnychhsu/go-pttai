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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/ailabstw/go-pttai/common"
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/crypto"
	"github.com/ailabstw/go-pttai/p2p"
)

var genIV = func(iv []byte) error {
	_, err := io.ReadFull(rand.Reader, iv)
	return err
}

// encrypt / decrypt refers:
// https://gist.github.com/stupidbodo/601b68bfef3449d1b8d9

func (p *BasePtt) EncryptData(op OpType, data []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	keyBytes := crypto.FromECDSA(key)
	marshaled := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(marshaled[:4], uint32(op))

	copy(marshaled[4:], data)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	msg := aesPad(marshaled)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]

	err = genIV(iv)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], msg)

	return ciphertext, nil
}

func (p *BasePtt) DecryptData(ciphertext []byte, key *ecdsa.PrivateKey) (OpType, []byte, error) {
	keyBytes := crypto.FromECDSA(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return 0, nil, err
	}

	if (len(ciphertext) % aes.BlockSize) != 0 {
		return 0, nil, ErrInvalidData
	}

	iv := ciphertext[:aes.BlockSize]
	msg := ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	marshaled, err := aesUnpad(msg)
	if err != nil {
		return 0, nil, err
	}

	opBytes := marshaled[:4]
	op := OpType(binary.BigEndian.Uint32(opBytes))
	data := marshaled[4:]

	return op, data, nil
}

func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

func removeBase64Padding(value string) string {
	return strings.TrimRight(value, "=")
}

func aesPad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func aesUnpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, ErrInvalidData
	}

	return src[:(length - unpadding)], nil
}

/*
MarshalData marshals the encrypted data based on ptt-protocol.
	hash: entity-hash
	enc: encrypted-data
The purpose is to have checksum to ensure that the data is not randomly-modified (preventing machine-error)
*/
func (p *BasePtt) MarshalData(code CodeType, hash *common.Address, encData []byte) (*PttData, error) {
	// 2. forms pttEvent
	ev := &PttEventData{
		Code:    code,
		Hash:    hash[:],
		EncData: encData,
	}

	// ptt-event signed
	evWithSalt, checksum, err := p.ChecksumPttEventData(ev)
	if err != nil {
		return nil, err
	}

	return &PttData{
		Code:       code,
		Hash:       hash[:],
		EvWithSalt: evWithSalt,
		Checksum:   checksum,
		Relay:      0,
	}, nil
}

/*
ChecksumPttEventData do checksum on the ev

Return: bytesWithSalt, checksum, error
*/
func (p *BasePtt) ChecksumPttEventData(ev *PttEventData) ([]byte, []byte, error) {
	evBytes, err := json.Marshal(ev)
	if err != nil {
		return nil, nil, err
	}

	return p.ChecksumData(evBytes)
}

/*
ChecksumData do checksum on the bytes

Return: bytesWithSalt, checksum, error
*/
func (p *BasePtt) ChecksumData(bytes []byte) ([]byte, []byte, error) {
	salt, err := types.NewSalt()
	if err != nil {
		return nil, nil, err
	}

	bytesWithSalt, err := common.Concat([][]byte{bytes, salt[:]})
	if err != nil {
		return nil, nil, err
	}
	hash := crypto.Keccak256(bytesWithSalt)

	return bytesWithSalt, hash, nil
}

/*
PttUnmarshalData unmarshal the pttData to the original data
*/
func (p *BasePtt) UnmarshalData(pttData *PttData) (CodeType, *common.Address, []byte, error) {
	ev, err := p.VerifyChecksumEventData(pttData)
	if err != nil {
		return CodeTypeDummy, nil, nil, err
	}

	hashAddr := &common.Address{}
	copy(hashAddr[:], ev.Hash[:])

	return ev.Code, hashAddr, ev.EncData, nil
}

func (p *BasePtt) VerifyChecksumEventData(pttData *PttData) (*PttEventData, error) {
	evWithSalt, checksum := pttData.EvWithSalt, pttData.Checksum
	err := p.VerifyChecksumData(evWithSalt, checksum)
	if err != nil {
		return nil, err

	}

	evBytes := evWithSalt[:len(evWithSalt)-types.SizeSalt]

	ev := &PttEventData{}
	err = json.Unmarshal(evBytes, ev)
	if err != nil {
		return nil, err
	}

	return ev, nil

}

func (p *BasePtt) VerifyChecksumData(bytesWithSalt []byte, checksum []byte) error {
	hash := crypto.Keccak256(bytesWithSalt)

	isGood := reflect.DeepEqual(hash, checksum)
	if !isGood {
		return ErrInvalidData
	}
	return nil
}

/*
SignEventData Signs the PttEventData
*/
func (p *BasePtt) SignEventData(ev *PttEventData, key *ecdsa.PrivateKey) ([]byte, []byte, []byte, error) {
	evBytes, err := json.Marshal(ev)
	if err != nil {
		return nil, nil, nil, err
	}

	bytesWithSalt, _, sig, pubBytes, err := SignData(evBytes, key)
	return bytesWithSalt, sig, pubBytes, err
}

func (p *BasePtt) VerifyEventData(evWithSalt []byte, sig []byte, keyBytes []byte) (*PttEventData, error) {
	err := VerifyData(evWithSalt, sig, keyBytes)
	if err != nil {
		return nil, err

	}

	evBytes := evWithSalt[:len(evWithSalt)-types.SizeSalt]

	ev := &PttEventData{}
	err = json.Unmarshal(evBytes, ev)
	if err != nil {
		return nil, err
	}

	return ev, nil

}

func (p *BasePtt) Server() *p2p.Server {
	return p.server
}

func (p *BasePtt) NoMorePeers() chan struct{} {
	return p.noMorePeers
}
