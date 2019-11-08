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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/commands/commandstest"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts/promptstest"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/keystoretest"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_accountListCmd(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		ks keystore.Store
	}
	tests := []struct {
		name        string
		args        args
		cmdArgs     []string
		cmdFlags    map[string]string
		wantOutput  string
		wantExecErr bool
	}{
		{
			"success",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					m.EXPECT().GetAddresses("ethereum", "mainnet").Return([][]byte{
						testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
						testutil.MustHexDecodeString("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
					}, nil)
					return m
				}(),
			},
			nil,
			map[string]string{
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761\n4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2\n",
			false,
		},
		{
			"err-store",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					m.EXPECT().GetAddresses("ethereum", "mainnet").Return([][]byte{}, errors.Errorf("failed"))
					return m
				}(),
			},
			nil,
			map[string]string{
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"Error: could not get addresses: failed",
			true,
		},
		{
			"err-empty-protocol",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					return m
				}(),
			},
			nil,
			map[string]string{
				"protocol": "",
				"network":  "mainnet",
			},
			"Error: `--protocol` must be specified to return address list",
			true,
		},
		{
			"err-empty-network",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					return m
				}(),
			},
			nil,
			map[string]string{
				"protocol": "ethereum",
				"network":  "",
			},
			"Error: `--network` must be specified to return address list",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := accountListCmd(tt.args.ks)
			if !assert.NotNil(got) {
				t.Error("accountListCmd() is nil")
			}
			_, out, err := commandstest.ExecuteCommandC(got, tt.cmdArgs, tt.cmdFlags)
			if (err != nil) != tt.wantExecErr {
				t.Errorf("configChainEthereumNetwork().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
				return
			}
			if !commandstest.AssertCommandOutput(t, got, err, out, tt.wantOutput) {
				t.Errorf("configChainEthereumNetwork().Execute().out != %v", tt.wantOutput)
			}
		})
	}
}

func Test_accountAddCmd(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		ks               keystore.Store
		passphrasePrompt func(suppliedSecret string, prePromptNote string, promptLabel string, allowEmpty bool, confirmPrompt bool) (string, error)
		privateKeyPrompt func(suppliedSecret string, prePromptNote string, promptLabel string, allowEmpty bool, confirmPrompt bool) (string, error)
	}
	tests := []struct {
		name        string
		args        args
		cmdArgs     []string
		cmdFlags    map[string]string
		wantOutput  string
		wantExecErr bool
	}{
		{
			"success",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					pk, _ := multikey.PrivateKeyFromBytes(secp256k1test.SofiaPrivateKey.Kind(), secp256k1test.SofiaPrivateKey.Bytes())
					m.EXPECT().Store(pk, gomock.Any()).Return(secp256k1test.SofiaPublicKey, nil)
					return m
				}(),
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, "01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F", nil),
			},
			nil,
			map[string]string{
				"key-type": crypto.SECP256K1,
			},
			"\x1b[32mPrivate key added\n\x1b[39mPublic key=0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006\n",
			false,
		},
		{
			"err-keystore",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					pk, _ := multikey.PrivateKeyFromBytes(secp256k1test.SofiaPrivateKey.Kind(), secp256k1test.SofiaPrivateKey.Bytes())
					m.EXPECT().Store(pk, gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m
				}(),
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, "01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F", nil),
			},
			nil,
			map[string]string{
				"key-type": crypto.SECP256K1,
			},
			"Error: key could not be stored: failed",
			true,
		},
		{
			"err-passphrase",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					return m
				}(),
				promptstest.MockRequiredSecret(t, "", errors.Errorf("failed")),
				promptstest.MockRequiredSecret(t, "01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F", nil),
			},
			nil,
			map[string]string{
				"key-type": crypto.SECP256K1,
			},
			"Error: could not get `passphrase`: failed",
			true,
		},
		{
			"err-private-key-invalid",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					return m
				}(),
				nil,
				promptstest.MockRequiredSecret(t, "01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A41594F", nil),
			},
			nil,
			map[string]string{
				"key-type": crypto.SECP256K1,
			},
			"Error: `private-key` could not be decoded: encoding/hex: odd length hex string",
			true,
		},
		{
			"err-private-key",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					return m
				}(),
				nil,
				promptstest.MockRequiredSecret(t, "", errors.Errorf("failed")),
			},
			nil,
			map[string]string{
				"key-type": crypto.SECP256K1,
			},
			"Error: could not get private key: failed",
			true,
		},
		{
			"err-key-type",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					return m
				}(),
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, "01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F", nil),
			},
			nil,
			map[string]string{
				"key-type": "invalid",
			},
			"Error: `private-key` could not be created from bytes: unsupported key type: \"invalid\"",
			true,
		},
		{
			"err-empty-key-type",
			args{
				func() keystore.Store {
					m := keystoretest.NewMockStore(mockCtrl)
					return m
				}(),
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, "01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F", nil),
			},
			nil,
			map[string]string{
				"key-type": "",
			},
			"Error: `key-type` must be specified",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := accountAddCmd(tt.args.ks, tt.args.passphrasePrompt, tt.args.privateKeyPrompt)
			if !assert.NotNil(got) {
				t.Error("accountListCmd() is nil")
			}
			_, out, err := commandstest.ExecuteCommandC(got, tt.cmdArgs, tt.cmdFlags)
			if (err != nil) != tt.wantExecErr {
				t.Errorf("configChainEthereumNetwork().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
				return
			}
			if !commandstest.AssertCommandOutput(t, got, err, out, tt.wantOutput) {
				t.Errorf("configChainEthereumNetwork().Execute().out != %v", tt.wantOutput)
			}
		})
	}
}

func Test_accountCmd(t *testing.T) {
	type args struct {
		config *settings.Root
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantNil     bool
		cmdArgs     []string
		cmdFlags    map[string]string
		wantExecErr bool
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
			false,
			false,
			nil,
			map[string]string{},
			false,
		},
		{
			"err-keystore",
			args{
				func() *settings.Root {
					v := viper.New()
					v.Set("keystore.kind", "invalid")
					config := settings.FromStore(v)
					return config
				}(),
			},
			true,
			true,
			nil,
			map[string]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := accountCmd(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("getKeyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Error("accountListCmd() is nil")
			}
			if !tt.wantNil {
				_, _, err = commandstest.ExecuteCommandC(got, tt.cmdArgs, tt.cmdFlags)
				if (err != nil) != tt.wantExecErr {
					t.Errorf("configChainEthereumNetwork().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
					return
				}
				if (err != nil) != tt.wantErr {
					t.Errorf("accountCmd() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
