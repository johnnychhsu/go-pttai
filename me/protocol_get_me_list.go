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
	"reflect"

	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/content"
)

func (spm *ServiceProtocolManager) GetMeList(contentBackend *content.Backend, startingID *types.PttID, limit int) (*MyInfo, []*MyInfo, error) {
	iter, err := getMeIter(startingID)
	if err != nil {
		return nil, nil, err
	}
	defer iter.Release()

	meList := make([]*MyInfo, 0)
	i := 0

	var myInfo *MyInfo
	for iter.Next() {
		if limit > 0 && i >= limit {
			break
		}

		v := iter.Value()

		eachMe := &MyInfo{}
		err = eachMe.Unmarshal(v)
		if err != nil {
			continue
		}

		if reflect.DeepEqual(MyID, eachMe.ID) {
			myInfo = eachMe
		}

		meList = append(meList, eachMe)

		i++
	}

	return myInfo, meList, nil
}
