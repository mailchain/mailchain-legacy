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

	"github.com/mailchain/mailchain/cmd/mailchain/internal/commands/prompts"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/s3store"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSentStore_setS3(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	tests := []struct {
		name         string
		fields       fields
		wantErr      bool
		wantSettings map[string]interface{}
	}{
		{
			"success-empty",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return label + "-value", nil
				},
			},
			false,
			map[string]interface{}{
				"stores": map[string]interface{}{
					"s3": map[string]interface{}{
						"access-key-id":     "access-key-id-value",
						"bucket":            "bucket-value",
						"region":            "region-value",
						"secret-access-key": "secret-access-key-value",
					}}},
		},
		{
			"err-empty-bucket",
			fields{
				viper.New(),
				func(label string) (string, error) {
					if label == "bucket" {
						return "", errors.Errorf("prompt failed")
					}
					return label + "-value", nil
				},
			},
			true,
			map[string]interface{}{},
		},
		{
			"err-empty-region",
			fields{
				viper.New(),
				func(label string) (string, error) {
					if label == "region" {
						return "", errors.Errorf("prompt failed")
					}
					return label + "-value", nil
				},
			},
			true,
			map[string]interface{}{},
		},
		{
			"err-empty-access-key-id",
			fields{
				viper.New(),
				func(label string) (string, error) {
					if label == "access-key-id" {
						return "", errors.Errorf("prompt failed")
					}
					return label + "-value", nil
				},
			},
			true,
			map[string]interface{}{},
		},
		{
			"err-empty-secret-access-key",
			fields{
				viper.New(),
				func(label string) (string, error) {
					if label == "secret-access-key" {
						return "", errors.Errorf("prompt failed")
					}
					return label + "-value", nil
				},
			},
			true,
			map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SentStore{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			if err := s.setS3(); (err != nil) != tt.wantErr {
				t.Errorf("SentStores.setS3() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.Equal(tt.wantSettings, tt.fields.viper.AllSettings()) {
				t.Errorf("settings = %v, wantSettings %v", tt.fields.viper.AllSettings(), tt.wantSettings)
			}
		})
	}
}

func TestSentStore_Set(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	type args struct {
		sentType string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantSettings map[string]interface{}
	}{
		{
			"success-default",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return label + "-value", nil
				},
			},
			args{
				"mailchain",
			},
			false,
			map[string]interface{}{
				"storage": map[string]interface{}{
					"sent": "mailchain"},
			},
		},
		{
			"success-s3",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return label + "-value", nil
				},
			},
			args{
				"s3",
			},
			false,
			map[string]interface{}{
				"storage": map[string]interface{}{
					"sent": "s3"},
				"stores": map[string]interface{}{
					"s3": map[string]interface{}{
						"access-key-id":     "access-key-id-value",
						"bucket":            "bucket-value",
						"region":            "region-value",
						"secret-access-key": "secret-access-key-value",
					},
				}},
		},
		{
			"err-unknown",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return label + "-value", nil
				},
			},
			args{
				"invalid",
			},
			true,
			map[string]interface{}{
				"storage": map[string]interface{}{
					"sent": "invalid"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SentStore{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			if err := s.Set(tt.args.sentType); (err != nil) != tt.wantErr {
				t.Errorf("SentStores.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.Equal(tt.wantSettings, tt.fields.viper.AllSettings()) {
				t.Errorf("settings = %v, wantSettings %v", tt.fields.viper.AllSettings(), tt.wantSettings)
			}
		})
	}
}

func TestSentStore_Get(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    stores.Sent
		wantErr bool
	}{
		{
			"invalid",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.sent", "invalid")
					return v
				}(),
				nil,
			},
			nil,
			true,
		},
		{
			"empty",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.sent", "")
					return v
				}(),
				nil,
			},
			stores.SentStore{},
			false,
		},
		{
			"mailchain",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.sent", "mailchain")
					return v
				}(),
				nil,
			},
			stores.SentStore{},
			false,
		},
		{
			"s3",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.sent", "s3")
					v.Set("stores.s3.region", "us-east-1")
					v.Set("stores.s3.bucket", "bucket")
					return v
				}(),
				nil,
			},
			&s3store.Sent{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SentStore{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			got, err := s.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("SentStore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.want, got) {
				t.Errorf("SentStore.Get() = %T, want %v]T", got, tt.want)
			}
		})
	}
}

func TestDefaultSentStore(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want SentStore
	}{
		{
			"success",
			SentStore{
				viper:         viper.GetViper(),
				requiredInput: prompts.RequiredInput,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultSentStore()
			if !assert.EqualValues(tt.want.viper, got.viper) {
				t.Errorf("DefaultSentStore().viper = %v, want %v", got.viper, tt.want.viper)
			}
		})
	}
}
