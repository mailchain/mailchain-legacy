// Copyright 2022 Mailchain Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocols

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkNames(t *testing.T) {
	type args struct {
		chain string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"algorand",
			args{
				"algorand",
			},
			[]string{"mainnet", "betanet", "testnet"},
		},
		{
			"ethereum",
			args{
				"ethereum",
			},
			[]string{"goerli", "kovan", "mainnet", "rinkeby", "ropsten"},
		},
		{
			"substrate",
			args{
				"substrate",
			},
			[]string{"edgeware-mainnet", "edgeware-beresheet", "edgeware-local"},
		},
		{
			"unknown",
			args{
				"unknown",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NetworkNames(tt.args.chain); !assert.Equal(t, tt.want, got) {
				t.Errorf("NetworkNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			"success",
			[]string{"algorand", "ethereum", "substrate"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := All(); !assert.Equal(t, tt.want, got) {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}
