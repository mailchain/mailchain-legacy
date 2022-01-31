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

package nameservice

import (
	errs "errors"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_WrapError(t *testing.T) {
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
			if (err != nil) && (tt.wantErrMessage != "") && !assert.EqualError(t, err, tt.wantErrMessage) {
				t.Errorf("wrapError() errorMessage = %v, wantErrMessage %v", err, tt.wantErrMessage)
			}
		})
	}
}

func TestErrorToRFC1035Status(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"err-format",
			args{ErrFormat},
			1,
		},
		{
			"err-serv-fail",
			args{ErrServFail},
			2,
		},
		{
			"err-nx-domain",
			args{ErrNXDomain},
			3,
		},
		{
			"err-not-imp",
			args{ErrNotImp},
			4,
		},
		{
			"err-refused",
			args{ErrRefused},
			5,
		},
		{
			"nil",
			args{nil},
			0,
		},
		{
			"err-other",
			args{errors.Errorf("error")},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorToRFC1035Status(tt.args.err); got != tt.want {
				t.Errorf("ErrorToRFC1035Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRFC1035StatusToError(t *testing.T) {
	type args struct {
		status int
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			"err-format",
			args{1},
			ErrFormat,
		},
		{
			"err-serv-fail",
			args{2},
			ErrServFail,
		},
		{
			"err-nx-domain",
			args{3},
			ErrNXDomain,
		},
		{
			"err-not-imp",
			args{4},
			ErrNotImp,
		},
		{
			"err-refused",
			args{5},
			ErrRefused,
		},
		{
			"nil",
			args{0},
			nil,
		},
		{
			"err-other",
			args{-1},
			errs.New("unknown RFC1035 status: -1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RFC1035StatusToError(tt.args.status); !assert.Equal(t, err, tt.wantErr) {
				t.Errorf("RFC1035StatusToError() errorMessage = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
