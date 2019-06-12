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
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts/promptstest"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestKeystore_selectKeystore(t *testing.T) {
	type fields struct {
		setter             config.KeystoreSetter
		viper              *viper.Viper
		selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
	}
	type args struct {
		keystoreType string
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
					v.Set("storage.keys", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"nacl-filestore"}, "already-set", true, nil),
			},
			args{
				mailchain.RequiresValue,
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
					v.Set("storage.keys", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"nacl-filestore"}, "new-value", false, nil),
			},
			args{
				mailchain.RequiresValue,
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
					v.Set("storage.keys", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"nacl-filestore"}, "", false, errors.Errorf("failed to select")),
			},
			args{
				mailchain.RequiresValue,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := Keystore{
				setter:             tt.fields.setter,
				viper:              tt.fields.viper,
				selectItemSkipable: tt.fields.selectItemSkipable,
			}
			got, err := k.selectKeystore(tt.args.keystoreType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keystore.selectKeystore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Keystore.selectKeystore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeystore_Select(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		setter             config.KeystoreSetter
		viper              *viper.Viper
		selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
	}
	type args struct {
		cmd          *cobra.Command
		keystoreType string
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
				func() config.KeystoreSetter {
					setter := configtest.NewMockKeystoreSetter(mockCtrl)
					setter.EXPECT().Set("new-keystore-type-value", "cmd-keystore-path-value").Return(nil)
					return setter
				}(),
				nil,
				nil,
			},
			args{
				func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("keystore-path", "cmd-keystore-path-value", "")

					return cmd
				}(),
				"new-keystore-type-value",
			},
			"new-keystore-type-value",
			false,
		},
		{
			"err-select-failed",
			fields{
				func() config.KeystoreSetter {
					setter := configtest.NewMockKeystoreSetter(mockCtrl)
					return setter
				}(),
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"nacl-filestore"}, "", true, errors.Errorf("failed to skip")),
			},
			args{
				func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("keystore-path", "cmd-keystore-path-value", "")

					return cmd
				}(),
				mailchain.RequiresValue,
			},
			"",
			true,
		},
		{
			"err-setter-failed",
			fields{
				func() config.KeystoreSetter {
					setter := configtest.NewMockKeystoreSetter(mockCtrl)
					setter.EXPECT().Set("new-setting", "cmd-keystore-path-value").Return(errors.Errorf("set failed"))
					return setter
				}(),
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{"nacl-filestore"}, "new-setting", false, nil),
			},
			args{
				func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("keystore-path", "cmd-keystore-path-value", "")

					return cmd
				}(),
				mailchain.RequiresValue,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := Keystore{
				setter:             tt.fields.setter,
				viper:              tt.fields.viper,
				selectItemSkipable: tt.fields.selectItemSkipable,
			}
			got, err := k.Select(tt.args.cmd, tt.args.keystoreType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keystore.Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Keystore.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}
