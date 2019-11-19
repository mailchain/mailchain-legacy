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

package nameservice

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"
)

func Test_WrapError(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		err error
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantErrMessage string
	}{
		{
			"no-resolver",
			args{
				errors.Errorf("%s: %s", "error", noResolverErrorMsg),
			},
			true,
			ErrNXDomain.Error(),
		},
		{
			"no-resolution",
			args{
				errors.Errorf("%s: %s", "error", noResolutionErrorMsg),
			},
			true,
			ErrNXDomain.Error(),
		},
		{
			"unregistered-name",
			args{
				errors.Errorf("%s: %s", "error", unregisteredNameErrorMsg),
			},
			true,
			ErrNXDomain.Error(),
		},
		{
			"could-not-parse-address",
			args{
				errors.Errorf("%s: %s", "error", couldNotParseAddressErrorMsg),
			},
			true,
			ErrFormat.Error(),
		},
		{
			"unknown error",
			args{
				errors.Errorf("%s: %s", "error", "original"),
			},
			true,
			"error: original",
		},
		{
			"nil error",
			args{
				nil,
			},
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapError(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("wrapError() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && (tt.wantErrMessage != "") && !assert.EqualError(err, tt.wantErrMessage) {
				t.Errorf("wrapError() errorMessage = %v, wantErrMessage %v", err, tt.wantErrMessage)
			}
		})
	}
}

func TestIsRFC1035Error(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"RFC1035Error-ErrFormat",
			args{
				ErrFormat,
			},
			true,
		},
		{
			"RFC1035Error-ErrServFail",
			args{
				ErrServFail,
			},
			true,
		},
		{
			"RFC1035Error-ErrNXDomain",
			args{
				ErrNXDomain,
			},
			true,
		},
		{
			"RFC1035Error-ErrNotImp",
			args{
				ErrNotImp,
			},
			true,
		},
		{
			"RFC1035Error-ErrRefused",
			args{
				ErrRefused,
			},
			true,
		},
		{
			"other",
			args{
				errors.Errorf("other"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRFC1035Error(tt.args.err); got != tt.want {
				t.Errorf("IsInvalidAddressError() = %v, want %v", got, tt.want)
			}
		})
	}
}
