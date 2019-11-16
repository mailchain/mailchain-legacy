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

package commands

import (
	"strings"
	"testing"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_settingsCmd(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		config *settings.Root
	}
	tests := []struct {
		name             string
		args             args
		wantCommandNames []string
	}{
		{
			"success",
			args{
				func() *settings.Root {
					v := viper.New()
					config := settings.FromStore(v)
					return config
				}(),
			},
			[]string{
				"view",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := settingsCmd(tt.args.config)
			subCommandNames := []string{}
			for _, x := range got.Commands() {
				subCommandNames = append(subCommandNames, x.Name())
			}
			if !assert.Equal(tt.wantCommandNames, subCommandNames) {
				t.Errorf("settingsCmd().Commands = %v, wantCommandNames %v", subCommandNames, strings.Join(tt.wantCommandNames, ","))
			}
		})
	}
}

func Test_settingsViewAll(t *testing.T) {
	type args struct {
		config *settings.Root
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{
				func() *settings.Root {
					v := viper.New()
					v.Set("server.port", 99999)
					config := settings.FromStore(v)
					return config
				}(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := settingsViewAll(tt.args.config)
			err := got.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("settingsViewAll().Execute error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
