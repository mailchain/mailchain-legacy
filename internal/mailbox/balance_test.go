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

package mailbox

import (
	"testing"

	"github.com/pkg/errors"
)

func TestIsBalanceNotSupportedError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"nil",
			args{
				nil,
			},
			false,
		},
		{
			"alt-err",
			args{
				errors.New("different error"),
			},
			false,
		},
		{
			"simple",
			args{
				errors.New("network not supported"),
			},
			true,
		},
		{
			"cause",
			args{
				errors.WithMessage(errors.New("network not supported"), "some other error"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNetworkNotSupportedError(tt.args.err); got != tt.want {
				t.Errorf("IsNetworkNotSupportedError() = %v, want %v", got, tt.want)
			}
		})
	}
}
