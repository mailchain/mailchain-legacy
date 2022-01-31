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

package mail

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHeaders(t *testing.T) {
	type args struct {
		date        time.Time
		from        Address
		to          Address
		replyTo     *Address
		subject     string
		contentType string
	}
	tests := []struct {
		name string
		args args
		want *Headers
	}{
		{
			"simple",
			args{
				time.Date(2001, 01, 02, 03, 04, 5, 6, time.UTC),
				Address{ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				Address{ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb", DisplayName: "", FullAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum"},
				nil,
				"Hello World",
				"text/plain; charset=\"UTF-8\"",
			},
			&Headers{
				Date:        time.Date(2001, 01, 02, 03, 04, 5, 6, time.UTC),
				From:        Address{ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				To:          Address{ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb", DisplayName: "", FullAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum"},
				Subject:     "Hello World",
				ContentType: "text/plain; charset=\"UTF-8\"",
			},
		},
		{
			"html",
			args{
				time.Date(2001, 01, 02, 03, 04, 5, 6, time.UTC),
				Address{ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				Address{ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb", DisplayName: "", FullAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum"},
				nil,
				"Hello World",
				"text/html; charset=\"UTF-8\"",
			},
			&Headers{
				Date:        time.Date(2001, 01, 02, 03, 04, 5, 6, time.UTC),
				From:        Address{ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				To:          Address{ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb", DisplayName: "", FullAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum"},
				Subject:     "Hello World",
				ContentType: "text/html; charset=\"UTF-8\"",
			},
		},
		{
			"reply-to",
			args{
				time.Date(2001, 01, 02, 03, 04, 5, 6, time.UTC),
				Address{ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				Address{ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb", DisplayName: "", FullAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum"},
				&Address{ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
				"Hello World",
				"text/plain; charset=\"UTF-8\"",
			},
			&Headers{
				Date:        time.Date(2001, 01, 02, 03, 04, 5, 6, time.UTC),
				From:        Address{ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				To:          Address{ChainAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb", DisplayName: "", FullAddress: "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum"},
				ReplyTo:     &Address{ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
				Subject:     "Hello World",
				ContentType: "text/plain; charset=\"UTF-8\"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHeaders(tt.args.date, tt.args.from, tt.args.to, tt.args.replyTo, tt.args.subject, tt.args.contentType)
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("NewHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
