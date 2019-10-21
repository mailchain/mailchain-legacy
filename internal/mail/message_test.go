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
	"time"
)

func TestNewMessage(t *testing.T) {
	type args struct {
		date    time.Time
		from    Address
		to      Address
		replyTo *Address
		subject string
		body    []byte
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			args{
				time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
				Address{
					FullAddress:  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum",
					ChainAddress: "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				},
				Address{
					FullAddress:  "0x4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum",
					ChainAddress: "0x4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2",
				},
				nil,
				"test subject",
				[]byte("test body"),
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMessage(tt.args.date, tt.args.from, tt.args.to, tt.args.replyTo, tt.args.subject, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("NewMessage() nil = %v, wantNil %v", err, tt.wantNil)
				return
			}
		})
	}
}

func Test_detectContentType(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"content-type-and-encoding",
			args{
				[]byte("this is plain text message"),
			},
			"text/plain; charset=\"UTF-8\"",
		},
		{
			"content-type-and-encoding-html",
			args{
				[]byte("<h1>this is HTML text message</h1>"),
			},
			"text/html; charset=\"UTF-8\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContentType(tt.args.body)
			if got != tt.want {
				t.Errorf("detectContentType() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
