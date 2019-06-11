// Copyright 2019 Finobo
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

package chains

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/internal/encoding"
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
			"ethereum",
			args{
				"ethereum",
			},
			encoding.EthereumNetworks(),
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
			if got := NetworkNames(tt.args.chain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworkNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
