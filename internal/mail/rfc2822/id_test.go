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

package rfc2822

import (
	nm "net/mail"
	"testing"

	"github.com/mailchain/mailchain/internal/mail"
	"github.com/stretchr/testify/assert"
)

func Test_parseID(t *testing.T) {
	type args struct {
		h nm.Header
	}
	tests := []struct {
		name    string
		args    args
		want    mail.ID
		wantErr bool
	}{
		{
			"success",
			args{
				nm.Header{
					"Message-Id": []string{"47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471@mailchain"},
				},
			},
			mail.ID{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			false,
		},
		{
			"err-decode",
			args{
				nm.Header{
					"Message-Id": []string{"47eca01e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471@mailchain"},
				},
			},
			nil,
			true,
		},
		{
			"err-no-suffix",
			args{
				nm.Header{
					"Message-Id": []string{"47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471@something"},
				},
			},
			nil,
			true,
		},
		{
			"err-empty-header",
			args{
				nm.Header{
					"Message-Id": []string{},
				},
			},
			nil,
			true,
		},
		{
			"err-no-header",
			args{
				nm.Header{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseID(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseID() = %v, want %v", got, tt.want)
			}
		})
	}
}
