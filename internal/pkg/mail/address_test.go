// Copyright (c) 2019 Finobo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
