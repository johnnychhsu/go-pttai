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
	"github.com/ailabstw/go-pttai/p2p/discover"
	"github.com/ailabstw/go-pttai/pttdb"
)

/**********
 * Me
 **********/

func (p *BasePtt) MyNodeID() *discover.NodeID {
	return p.myNodeID
}

func (p *BasePtt) SetMyEntity(myEntity PttMyEntity) error {
	var err error
	p.myEntity = myEntity

	p.meOplogMerkle, err = NewMerkle(DBMeOplogPrefix, DBMeMerkleOplogPrefix, myEntity.GetID(), dbOplog)
	if err != nil {
		return err
	}

	p.masterOplogMerkle, err = NewMerkle(DBMasterOplogPrefix, DBMasterMerkleOplogPrefix, myEntity.GetID(), dbOplog)
	if err != nil {
		return err
	}

	return nil
}

func (p *BasePtt) GetMyEntity() MyEntity {
	return p.myEntity
}

func (p *BasePtt) SignKey() *KeyInfo {
	return p.myEntity.SignKey()
}

func (p *BasePtt) MeOplogMerkle() *Merkle {
	return p.meOplogMerkle
}

func (p *BasePtt) MasterOplogMerkle() *Merkle {
	return p.masterOplogMerkle
}

func (p *BasePtt) DBOplog() *pttdb.LDBBatch {
	return dbOplog
}
