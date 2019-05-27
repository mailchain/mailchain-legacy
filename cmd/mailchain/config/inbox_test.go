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

	"github.com/mailchain/mailchain/internal/stores"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetStateStore(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		vpr *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		want    stores.State
		wantErr bool
	}{
		// {
		// 	"success",
		// 	args{
		// 		func() *viper.Viper {
		// 			v := viper.New()
		// 			v.Set("storage.inbox", "leveldb")
		// 			v.Set("stores.leveldb.path", "./")
		// 			return v
		// 		}(),
		// 	},
		// 	nil,
		// 	false,
		// },
		{
			"error",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.state", "invalid")
					return v
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStateStore(tt.args.vpr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStateStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("GetStateStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
