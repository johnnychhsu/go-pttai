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

	"github.com/ailabstw/go-pttai/common/types"
)

/*
IdentifyPeerWithMyIDChallenge requests challenge to make sure the user-id (acker)
*/
func (p *BasePtt) IdentifyPeerWithMyIDChallenge(userID *types.PttID, peer *PttPeer) error {
	data, err := p.IdentifyPeer(userID, p.quitSync, peer)
	if err != nil {
		return err
	}

	return p.SendDataToPeer(CodeTypeIdentifyPeerWithMyIDChallenge, data, peer)
}

/*
HandleIdentifyPeerWithMyIDChallenge handles IdentifyPeerWithMyIDChallenge (requester)
*/
func (p *BasePtt) HandleIdentifyPeerWithMyIDChallenge(dataBytes []byte, peer *PttPeer) error {
	data := &IdentifyPeer{}
	err := json.Unmarshal(dataBytes, data)
	if err != nil {
		return err
	}

	return p.IdentifyPeerWithMyIDChallengeAck(data, peer)
}
