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

package config

import (
	"testing"

	"github.com/matryer/is"
	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

func Test_setEtherscan(t *testing.T) {
	is := is.New(t)
	type args struct {
		vpr           *viper.Viper
		requiredInput func(label string) (string, error)
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected map[string]interface{}
	}{
		{
			"set-api-key",
			args{
				viper.New(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			false,
			map[string]interface{}{"api-key": "api-key-value"},
		},
		{
			"error",
			args{
				viper.New(),
				func(label string) (string, error) {
					return "", errors.Errorf("failed")
				},
			},
			true,
			nil,
		},
		{
			"already-specified",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("clients.etherscan.api-key", "already-set")
					return v
				}(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			false,
			map[string]interface{}{"api-key": "already-set"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setEtherscan(tt.args.vpr, tt.args.requiredInput); (err != nil) != tt.wantErr {
				t.Errorf("setEtherscan() error = %v, wantErr %v", err, tt.wantErr)
			}

			is.Equal(tt.expected, tt.args.vpr.Get("clients.etherscan"))
		})
	}
}
