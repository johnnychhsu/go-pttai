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

package e2e

import (
	"testing"

	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/me"
	"github.com/stretchr/testify/assert"
	baloo "gopkg.in/h2non/baloo.v3"
)

func TestMeBasic(t *testing.T) {
	NNodes = 1

	var bodyString string
	assert := assert.New(t)

	setupTest(t)
	defer teardownTest(t)

	t0 := baloo.New("http://127.0.0.1:9450")

	// 1. get
	bodyString = `{"id": "testID", "method": "me_get", "params": []}`

	me0_1 := &me.BackendMyInfo{}

	testCore(t0, bodyString, me0_1, t, true)

	assert.Equal(types.StatusAlive, me0_1.Status)

	// 2. get total weight
	bodyString = `{"id": "testID", "method": "me_getTotalWeight", "params": []}`

	var totalWeigtht0_2 uint32
	testCore(t0, bodyString, &totalWeigtht0_2, t, true)

	assert.Equal(uint32(me.WeightDesktop), totalWeigtht0_2)

	// 3. getRawMe
	bodyString = `{"id": "testID", "method": "me_getRawMe", "params": []}`

	me0_3 := &me.MyInfo{}

	testCore(t0, bodyString, me0_3, t, true)

	assert.Equal(types.StatusAlive, me0_3.Status)
	assert.Equal(me0_1.ID, me0_3.ID)
	assert.Equal(me0_3.ID, me0_3.OwnerID)
}
