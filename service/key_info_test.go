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
	"crypto/ecdsa"
	"reflect"
	"testing"
)

func Test_deriveKeyPBKDF2(t *testing.T) {
	// setup test
	setupTest(t)
	defer teardownTest(t)

	// define test-structure
	type args struct {
		masterKey *ecdsa.PrivateKey
	}

	// prepare test-cases
	tests := []struct {
		name    string
		args    args
		want    *ecdsa.PrivateKey
		want1   *KeyExtraInfo
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args:    args{masterKey: tDefaultKey},
			want:    tDerivedKey,
			want1:   tDerivedExtraInfo,
			wantErr: false,
		},
	}

	// run test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := deriveKeyPBKDF2(tt.args.masterKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("deriveKeyPBKDF2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deriveKeyPBKDF2() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {

				t.Errorf("deriveKeyPBKDF2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	// teardown test
}
