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

// ServiceConstructor is the function signature of the constructors needed to be
// registered for service instantiation.
type ServiceConstructor func(ctx *ServiceContext) (PttService, error)

// PeerType
type PeerType int

const (
	PeerTypeErr PeerType = iota
	PeerTypeRemoved
	PeerTypeRandom
	PeerTypeMember
	PeerTypeImportant
	PeerTypeMe
)

// NodeType
type NodeType int

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeMobile
	NodeTypeDesktop
	NodeTypeServer
)
