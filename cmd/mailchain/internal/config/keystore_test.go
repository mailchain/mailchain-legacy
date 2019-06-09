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

	"github.com/mailchain/mailchain/internal/keystore/nacl"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestKeystore_Get(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		viper                    *viper.Viper
		requiredInputWithDefault func(label string, defaultValue string) (string, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    *nacl.FileStore
		wantErr bool
	}{
		{
			"success-nacl",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.keys", "nacl-filestore")
					return v
				}(),
				nil,
			},
			&nacl.FileStore{},
			false,
		},
		{
			"err-invalid",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.keys", "invalid")
					return v
				}(),
				nil,
			},
			nil,
			true,
		},
		{
			"err-empty",
			fields{
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := Keystore{
				viper:                    tt.fields.viper,
				requiredInputWithDefault: tt.fields.requiredInputWithDefault,
			}
			got, err := k.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Keystore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.want, got) {
				t.Errorf("Keystore.Get() = %T, want %v]T", got, tt.want)
			}
		})
	}
}
