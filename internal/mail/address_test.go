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
	nm "net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddress(t *testing.T) {

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
		{"invalid-address",
			args{"Charlotte <0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum", "ethereum", "ropsten"},
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
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("ParseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tryAddChainNetwork(t *testing.T) {
	type args struct {
		input   string
		chain   string
		network string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"has-@",
			args{"0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum", "ethereum", "ropsten"},
			"0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum",
			false,
		},
		{
			"has-chevrons",
			args{"<0x92d8f10248c6a3953cc3692a894655ad05d61efb>", "ethereum", "ropsten"},
			"",
			true,
		},
		{
			"missing-network",
			args{"0x92d8f10248c6a3953cc3692a894655ad05d61efb", "", "ethereum"},
			"",
			true,
		},
		{
			"missing-chain",
			args{"0x92d8f10248c6a3953cc3692a894655ad05d61efb", "ropsten", ""},
			"",
			true,
		},
		{
			"success",
			args{"0x92d8f10248c6a3953cc3692a894655ad05d61efb", "ethereum", "ropsten"},
			"0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tryAddChainNetwork(tt.args.input, tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("tryAddChainNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tryAddChainNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromAddress(t *testing.T) {
	type args struct {
		address *nm.Address
	}
	tests := []struct {
		name    string
		args    args
		want    *Address
		wantErr bool
	}{
		{
			"success",
			args{
				&nm.Address{
					Name:    "name",
					Address: "address@domain",
				},
			},
			&Address{DisplayName: "name", FullAddress: "address@domain", ChainAddress: "address"},
			false,
		},
		{
			"err-no-@",
			args{
				&nm.Address{
					Name:    "name",
					Address: "address",
				},
			},
			nil,
			true,
		},
		{
			"err-nil",
			args{
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("fromAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddress_String(t *testing.T) {
	type fields struct {
		DisplayName  string
		FullAddress  string
		ChainAddress string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"success",
			fields{
				"Charlotte",
				"0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum",
				"0x92d8f10248c6a3953cc3692a894655ad05d61efb",
			},
			"\"Charlotte\" <0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Address{
				DisplayName:  tt.fields.DisplayName,
				FullAddress:  tt.fields.FullAddress,
				ChainAddress: tt.fields.ChainAddress,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("Address.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
