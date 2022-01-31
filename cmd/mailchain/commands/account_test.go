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

package commands

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/commands/commandstest"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts/promptstest"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/keystoretest"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_accountListCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		produceKeystore func() (keystore.Store, error)
	}
	tests := []struct {
		name          string
		args          args
		cmdArgs       []string
		cmdFlags      map[string]string
		wantErrOutput string
		wantExecErr   bool
	}{
		{
			"query-ethereum-mainnet",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					m.EXPECT().GetAddresses("ethereum", "mainnet").Return(map[string]map[string][][]byte{
						"ethereum": {
							"mainnet": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
						},
					}, nil)
					return m, nil
				},
			},
			nil,
			map[string]string{
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"",
			false,
		},
		{
			"query-ethereum",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					m.EXPECT().GetAddresses("ethereum", "").Return(map[string]map[string][][]byte{
						"ethereum": {
							"mainnet": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
							"goerli": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
						},
					}, nil)
					return m, nil
				},
			},
			nil,
			map[string]string{
				"protocol": "ethereum",
			},
			"",
			false,
		},
		{
			"query-substrate",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					m.EXPECT().GetAddresses("substrate", "").Return(map[string]map[string][][]byte{
						"substrate": {
							"edgeware-beresheet": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
							"edgeware-mainnet": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
						},
					}, nil)
					return m, nil
				},
			},
			nil,
			map[string]string{
				"protocol": "substrate",
			},
			"",
			false,
		},
		{
			"query-all",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					m.EXPECT().GetAddresses("", "").Return(map[string]map[string][][]byte{
						"ethereum": {
							"mainnet": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
							"goerli": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
						},
						"substrate": {
							"edgeware-beresheet": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
							"edgeware-mainnet": {
								encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
								encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							},
						},
					}, nil)
					return m, nil
				},
			},
			nil,
			map[string]string{},
			"",
			false,
		},
		{
			"err-store",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					m.EXPECT().GetAddresses("ethereum", "mainnet").Return(nil, errors.Errorf("failed"))
					return m, nil
				},
			},
			nil,
			map[string]string{
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"Error: could not get addresses: failed",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := accountListCmd(tt.args.produceKeystore)
			if !assert.NotNil(t, got) {
				t.Error("accountListCmd() is nil")
			}
			_, out, err := commandstest.ExecuteCommandC(got, tt.cmdArgs, tt.cmdFlags)
			if (err != nil) != tt.wantExecErr {
				t.Errorf("configChainEthereumNetwork().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
				return
			}
			if !commandstest.AssertCommandJsonOutput(t, got, err, out, tt.wantErrOutput) {
				t.Errorf("configChainEthereumNetwork().Execute().out != %v", tt.wantErrOutput)
			}
		})
	}
}

func Test_accountAddCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		produceKeyStore  func() (keystore.Store, error)
		passphrasePrompt func(suppliedSecret string, prePromptNote string, promptLabel string, allowEmpty bool, confirmPrompt bool) (string, error)
		privateKeyPrompt func(suppliedSecret string, prePromptNote string, promptLabel string, allowEmpty bool, confirmPrompt bool) (string, error)
	}
	tests := []struct {
		name        string
		args        args
		cmdArgs     []string
		cmdFlags    map[string]string
		wantExecErr string
	}{
		{
			"success-hex",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					pk, _ := multikey.PrivateKeyFromBytes("secp256k1", secp256k1test.AlicePrivateKey.Bytes())
					m.EXPECT().Store("ethereum", "mainnet", pk, gomock.Any()).Return(secp256k1test.AlicePublicKey, nil)
					return m, nil
				},
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, encoding.EncodeHex(secp256k1test.AlicePrivateKey.Bytes()), nil),
			},
			nil,
			map[string]string{
				"key-type": crypto.KindSECP256K1,
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"",
		},
		{
			"success-mnemonic-algorand",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					pk, _ := multikey.PrivateKeyFromBytes("ed25519", ed25519test.AlicePrivateKey.Bytes())
					m.EXPECT().Store("algorand", "mainnet", pk, gomock.Any()).Return(ed25519test.AlicePublicKey, nil)
					return m, nil
				},
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, func() string {
					s, err := encoding.EncodeMnemonicAlgorand(ed25519test.AlicePrivateKey.Bytes()[:32])
					assert.NoError(t, err)
					return s
				}(), nil),
			},
			nil,
			map[string]string{
				"key-type":             crypto.KindED25519,
				"private-key-encoding": encoding.KindMnemonicAlgorand,
				"protocol":             "algorand",
				"network":              "mainnet",
			},
			"",
		},
		{
			"err-keystore",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					pk, _ := multikey.PrivateKeyFromBytes("secp256k1", secp256k1test.AlicePrivateKey.Bytes())
					m.EXPECT().Store("ethereum", "mainnet", pk, gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m, nil
				},
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, encoding.EncodeHex(secp256k1test.AlicePrivateKey.Bytes()), nil),
			},
			nil,
			map[string]string{
				"key-type": crypto.KindSECP256K1,
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"Error: key could not be stored: failed",
		},
		{
			"err-passphrase",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					return m, nil
				},
				promptstest.MockRequiredSecret(t, "", errors.Errorf("failed")),
				promptstest.MockRequiredSecret(t, encoding.EncodeHex(secp256k1test.AlicePrivateKey.Bytes()), nil),
			},
			nil,
			map[string]string{
				"key-type": crypto.KindSECP256K1,
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"Error: could not get `passphrase`: failed",
		},
		{
			"err-private-key-invalid",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					return m, nil
				},
				nil,
				promptstest.MockRequiredSecret(t, "01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A41594F", nil),
			},
			nil,
			map[string]string{
				"key-type": crypto.KindSECP256K1,
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"Error: `private-key` could not be decoded: encoding/hex: odd length hex string",
		},
		{
			"err-private-key",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					return m, nil
				},
				nil,
				promptstest.MockRequiredSecret(t, "", errors.Errorf("failed")),
			},
			nil,
			map[string]string{
				"key-type": crypto.KindSECP256K1,
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"Error: could not get private key: failed",
		},
		{
			"err-key-type",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					return m, nil
				},
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, encoding.EncodeHex(secp256k1test.AlicePrivateKey.Bytes()), nil),
			},
			nil,
			map[string]string{
				"key-type": "invalid",
				"protocol": "ethereum",
				"network":  "mainnet",
			},
			"Error: `private-key` could not be created from bytes: unsupported key type: \"invalid\"",
		},
		{
			"err-empty-key-type",
			args{
				func() (keystore.Store, error) {
					m := keystoretest.NewMockStore(mockCtrl)
					return m, nil
				},
				promptstest.MockRequiredSecret(t, "passphrase-secret", nil),
				promptstest.MockRequiredSecret(t, encoding.EncodeHex(secp256k1test.AlicePrivateKey.Bytes()), nil),
			},
			nil,
			map[string]string{
				// "key-type": "",
			},
			"Error: required flag(s) \"key-type\", \"network\" not set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := accountAddCmd(tt.args.produceKeyStore, tt.args.passphrasePrompt, tt.args.privateKeyPrompt)
			if !assert.NotNil(t, got) {
				t.Error("accountListCmd() is nil")
			}

			_, out, err := commandstest.ExecuteCommandC(got, tt.cmdArgs, tt.cmdFlags)

			if !commandstest.AssertCommandJsonOutput(t, got, err, out, tt.wantExecErr) {
				t.Errorf("Test_accountAddCmd().Execute().out != %v", tt.wantExecErr)
			}

			if tt.wantExecErr == "" && !assert.NoError(t, err) {
				t.Errorf("Test_accountAddCmd().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
			}

			// if tt.wantExecErr != "" && !assert.EqualError(t, err, tt.wantExecErr) {
			// 	t.Errorf("Test_accountAddCmd().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
			// 	return
			// }

			// if err == nil {
			// 	goldenResponse, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s.json", t.Name()))
			// 	if err != nil {
			// 		assert.FailNow(t, err.Error())
			// 	}
			// 	if !commandstest.AssertCommandJsonOutput(t, got, err, out, tt.wantExecErr) {
			// 		t.Errorf("configChainEthereumNetwork().Execute().out != %v", tt.wantExecErr)
			// 	}
			// 	if !assert.JSONEq(t, string(goldenResponse), out) {
			// 		t.Errorf("command returned unexpected response: got %v want %v",
			// 			out, goldenResponse)
			// 	}
			// }

			// if tt.wantOutput != "" && !commandstest.AssertCommandOutput(t, got, err, out, tt.wantOutput) {
			// 	t.Errorf("configChainEthereumNetwork().Execute().out != %v", tt.wantOutput)
			// }
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
