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

package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddress(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		input   string
		chain   string
		network string
	}
	tests := []struct {
		name    string
		args    args
		want    *Address
		wantErr bool
	}{
		{"full-address",
			args{"0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum", "", ""},
			&Address{
				FullAddress:  "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum",
				ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb",
			},
			false,
		}, {
			"full-address-with-display-name",
			args{"Charlotte <0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum>", "", ""},
			&Address{
				DisplayName:  "Charlotte",
				FullAddress:  "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum",
				ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb",
			},
			false,
		}, {
			"chain-only-address",
			args{"Charlotte <0x92d8f10248c6a3953cc3692a894655ad05d61efb>", "", ""},
			nil,
			true,
		},
		{
			"chain-only-address-with-chain-and-network",
			args{"0x92d8f10248c6a3953cc3692a894655ad05d61efb", "ethereum", "ropsten"},
			&Address{
				DisplayName:  "",
				ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb",
				FullAddress:  "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum",
			},
			false,
		},
		{"chain-only-address-with-display-name",
			args{"Charlotte <0x92d8f10248c6a3953cc3692a894655ad05d61efb>", "ethereum", "ropsten"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAddress(tt.args.input, tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(got, tt.want) {
				t.Errorf("ParseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
