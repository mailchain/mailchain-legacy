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

package setup

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/configtest"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/names"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts/promptstest"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func TestSentStorage_selectSentStorage(t *testing.T) {
	type fields struct {
		setter             config.SentStoreSetter
		viper              *viper.Viper
		selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
	}
	type args struct {
		sentStorageType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"success-already-set",
			fields{
				nil,
				nil,
				nil,
			},
			args{
				"value-already-set",
			},
			"value-already-set",
			false,
		},
		{
			"success-skipped",
			fields{
				nil,
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.sent", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"mailchain", "s3"}, "already-set", true, nil),
			},
			args{
				names.RequiresValue,
			},
			"",
			false,
		},
		{
			"success-not-skipped",
			fields{
				nil,
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.sent", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"mailchain", "s3"}, "new-value", false, nil),
			},
			args{
				names.RequiresValue,
			},
			"new-value",
			false,
		},
		{
			"err-skipped",
			fields{
				nil,
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.sent", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"mailchain", "s3"}, "", false, errors.Errorf("failed to select")),
			},
			args{
				names.RequiresValue,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SentStorage{
				setter:             tt.fields.setter,
				viper:              tt.fields.viper,
				selectItemSkipable: tt.fields.selectItemSkipable,
			}
			got, err := s.selectSentStorage(tt.args.sentStorageType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentStorage.selectSentStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SentStorage.selectSentStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSentStorage_Select(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		setter             config.SentStoreSetter
		viper              *viper.Viper
		selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
	}
	type args struct {
		sentStorageType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"success-set",
			fields{
				func() config.SentStoreSetter {
					setter := configtest.NewMockSentStoreSetter(mockCtrl)
					setter.EXPECT().Set("new-set-value").Return(nil)
					return setter
				}(),
				nil,
				nil,
			},
			args{
				"new-set-value",
			},
			"new-set-value",
			false,
		},
		{
			"success-skipped",
			fields{
				func() config.SentStoreSetter {
					setter := configtest.NewMockSentStoreSetter(mockCtrl)
					return setter
				}(),
				func() *viper.Viper {
					v := viper.New()
					v.Set("storage.sent", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"mailchain", "s3"}, "already-set", true, nil),
			},
			args{
				names.RequiresValue,
			},
			"",
			false,
		},
		{
			"err-select-failed",
			fields{
				func() config.SentStoreSetter {
					setter := configtest.NewMockSentStoreSetter(mockCtrl)
					return setter
				}(),
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"mailchain", "s3"}, "", true, errors.Errorf("failed to skip")),
			},
			args{
				names.RequiresValue,
			},
			"",
			true,
		},
		{
			"err-setter-failed",
			fields{
				func() config.SentStoreSetter {
					setter := configtest.NewMockSentStoreSetter(mockCtrl)
					setter.EXPECT().Set("new-setting").Return(errors.Errorf("failed to error"))
					return setter
				}(),
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"mailchain", "s3"}, "new-setting", false, nil),
			},
			args{
				names.RequiresValue,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SentStorage{
				setter:             tt.fields.setter,
				viper:              tt.fields.viper,
				selectItemSkipable: tt.fields.selectItemSkipable,
			}
			got, err := s.Select(tt.args.sentStorageType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentStorage.Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SentStorage.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}
