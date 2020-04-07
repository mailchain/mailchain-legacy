// Copyright 2020 Finobo
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

package substrate

import (
	"errors"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cryptotest"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/stretchr/testify/assert"
)

func TestSigner_Sign(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		privateKey crypto.PrivateKey
	}
	type args struct {
		opts signer.Options
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantSigner *types.Address
		wantMethod *types.Call
		wantErr    bool
	}{
		{
			"sr25519-charlotte",
			fields{
				sr25519test.CharlottePrivateKey,
			},
			args{
				SignerOptions{
					types.Extrinsic{
						Version: 0x4,
						Signature: types.ExtrinsicSignatureV4{
							Signer: types.Address{
								IsAccountID:    false,
								AsAccountID:    types.AccountID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsAccountIndex: false,
								AsAccountIndex: 0x0,
							},
							Signature: types.MultiSignature{
								IsEd25519: false,
								AsEd25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsSr25519: false,
								AsSr25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsEcdsa:   false,
								AsEcdsa:   types.Bytes(nil),
							},
							Era: types.ExtrinsicEra{
								IsImmortalEra: false,
								IsMortalEra:   false,
								AsMortalEra:   types.MortalEra{First: 0x0, Second: 0x0},
							},
							Nonce: 0x0,
							Tip:   0x0,
						},
						Method: types.Call{
							CallIndex: types.CallIndex{SectionIndex: 0x11, MethodIndex: 0x2},
							Args:      types.Args{0xff, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0x0, 0x2, 0xf4, 0x1, 0x0, 0x1c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
						},
					},
					types.SignatureOptions{
						Era: types.ExtrinsicEra{
							IsImmortalEra: false,
							IsMortalEra:   false,
							AsMortalEra:   types.MortalEra{First: 0x0, Second: 0x0},
						},
						Nonce:       0x7b,
						Tip:         0x1,
						SpecVersion: 0x0,
						GenesisHash: types.Hash{0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0x30, 0xd1, 0xd2},
						BlockHash:   types.Hash{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x30, 0x31, 0x32},
					},
				},
			},
			&types.Address{
				IsAccountID:    true,
				AsAccountID:    types.AccountID{0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30},
				IsAccountIndex: false,
				AsAccountIndex: 0x0,
			},
			&types.Call{
				CallIndex: types.CallIndex{
					SectionIndex: 0x11,
					MethodIndex:  0x2,
				},
				Args: types.Args{0xff, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0x0, 0x2, 0xf4, 0x1, 0x0, 0x1c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			},
			false,
		},
		{
			"err-nil-opts",
			fields{
				sr25519test.CharlottePrivateKey,
			},
			args{
				nil,
			},
			nil,
			nil,
			true,
		},
		{
			"err-sign",
			fields{
				sr25519test.CharlottePrivateKey,
			},
			args{
				&SignerOptions{
					types.Extrinsic{
						Version: 0x3,
						Signature: types.ExtrinsicSignatureV4{
							Signer: types.Address{
								IsAccountID:    false,
								AsAccountID:    types.AccountID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsAccountIndex: false,
								AsAccountIndex: 0x0,
							},
							Signature: types.MultiSignature{
								IsEd25519: false,
								AsEd25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsSr25519: false,
								AsSr25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsEcdsa:   false,
								AsEcdsa:   types.Bytes(nil),
							},
							Era: types.ExtrinsicEra{
								IsImmortalEra: false,
								IsMortalEra:   false,
								AsMortalEra:   types.MortalEra{First: 0x0, Second: 0x0},
							},
							Nonce: 0x0,
							Tip:   0x0,
						},
						Method: types.Call{
							CallIndex: types.CallIndex{SectionIndex: 0x11, MethodIndex: 0x2},
							Args:      types.Args{0xff, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0x0, 0x2, 0xf4, 0x1, 0x0, 0x1c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
						},
					},
					types.SignatureOptions{
						Era: types.ExtrinsicEra{
							IsImmortalEra: false,
							IsMortalEra:   false,
							AsMortalEra:   types.MortalEra{First: 0x0, Second: 0x0},
						},
						Nonce:       0x7b,
						Tip:         0x1,
						SpecVersion: 0x0,
						GenesisHash: types.Hash{0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0x30, 0xd1, 0xd2},
						BlockHash:   types.Hash{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x30, 0x31, 0x32},
					},
				},
			},
			nil,
			nil,
			true,
		},
		{
			"err-invalid-private-key",
			fields{
				func() crypto.PrivateKey {
					m := cryptotest.NewMockPrivateKey(mockCtrl)
					m.EXPECT().Sign(gomock.Any()).Return([]byte("signed-data"), nil)
					return m
				}(),
			},
			args{
				SignerOptions{
					types.Extrinsic{
						Version: 0x4,
						Signature: types.ExtrinsicSignatureV4{
							Signer: types.Address{
								IsAccountID:    false,
								AsAccountID:    types.AccountID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsAccountIndex: false,
								AsAccountIndex: 0x0,
							},
							Signature: types.MultiSignature{
								IsEd25519: false,
								AsEd25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsSr25519: false,
								AsSr25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsEcdsa:   false,
								AsEcdsa:   types.Bytes(nil),
							},
							Era: types.ExtrinsicEra{
								IsImmortalEra: false,
								IsMortalEra:   false,
								AsMortalEra:   types.MortalEra{First: 0x0, Second: 0x0},
							},
							Nonce: 0x0,
							Tip:   0x0,
						},
						Method: types.Call{
							CallIndex: types.CallIndex{SectionIndex: 0x11, MethodIndex: 0x2},
							Args:      types.Args{0xff, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0x0, 0x2, 0xf4, 0x1, 0x0, 0x1c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
						},
					},
					types.SignatureOptions{
						Era: types.ExtrinsicEra{
							IsImmortalEra: false,
							IsMortalEra:   false,
							AsMortalEra:   types.MortalEra{First: 0x0, Second: 0x0},
						},
						Nonce:       0x7b,
						Tip:         0x1,
						SpecVersion: 0x0,
						GenesisHash: types.Hash{0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0x30, 0xd1, 0xd2},
						BlockHash:   types.Hash{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x30, 0x31, 0x32},
					},
				},
			},
			nil,
			nil,
			true,
		},
		{
			"err-sign-failed",
			fields{
				func() crypto.PrivateKey {
					m := cryptotest.NewMockPrivateKey(mockCtrl)
					m.EXPECT().Sign(gomock.Any()).Return([]byte{}, errors.New("error"))
					return m
				}(),
			},
			args{
				SignerOptions{
					types.Extrinsic{
						Version: 0x4,
						Signature: types.ExtrinsicSignatureV4{
							Signer: types.Address{
								IsAccountID:    false,
								AsAccountID:    types.AccountID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsAccountIndex: false,
								AsAccountIndex: 0x0,
							},
							Signature: types.MultiSignature{
								IsEd25519: false,
								AsEd25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsSr25519: false,
								AsSr25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
								IsEcdsa:   false,
								AsEcdsa:   types.Bytes(nil),
							},
							Era: types.ExtrinsicEra{
								IsImmortalEra: false,
								IsMortalEra:   false,
								AsMortalEra:   types.MortalEra{First: 0x0, Second: 0x0},
							},
							Nonce: 0x0,
							Tip:   0x0,
						},
						Method: types.Call{
							CallIndex: types.CallIndex{SectionIndex: 0x11, MethodIndex: 0x2},
							Args:      types.Args{0xff, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0x0, 0x2, 0xf4, 0x1, 0x0, 0x1c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
						},
					},
					types.SignatureOptions{
						Era: types.ExtrinsicEra{
							IsImmortalEra: false,
							IsMortalEra:   false,
							AsMortalEra:   types.MortalEra{First: 0x0, Second: 0x0},
						},
						Nonce:       0x7b,
						Tip:         0x1,
						SpecVersion: 0x0,
						GenesisHash: types.Hash{0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0x30, 0xd1, 0xd2},
						BlockHash:   types.Hash{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x30, 0x31, 0x32},
					},
				},
			},
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Signer{
				privateKey: tt.fields.privateKey,
			}
			gotSignedTransactionInterface, err := e.Sign(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signer.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			gotSignedTransaction, ok := gotSignedTransactionInterface.(*types.Extrinsic)
			if !ok {
				assert.FailNow(t, "invalid return type")
			}
			opts, ok := tt.args.opts.(SignerOptions)
			if !ok {
				assert.FailNow(t, "invalid signature options")
			}
			data, err := e.prepareData(&opts.Extrinsic, &opts.SignatureOptions)
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			if !assert.Equal(t, tt.wantMethod, &gotSignedTransaction.Method) {
				t.Errorf("Signer.Sign().Method = %v, want %v", gotSignedTransaction.Method, tt.wantMethod)
			}
			var signature []byte
			switch tt.fields.privateKey.(type) {
			case *sr25519.PrivateKey:
				signature, _ = types.EncodeToBytes(gotSignedTransaction.Signature.Signature.AsSr25519)
			default:
				t.Error("unsupported key type")
			}

			verify := e.privateKey.PublicKey().Verify(data, signature)
			if !verify {
				t.Errorf("signature can not be verified by public key")
				return
			}

		})
	}
}

func TestNewSigner(t *testing.T) {
	type args struct {
		privateKey crypto.PrivateKey
	}
	tests := []struct {
		name string
		args args
		want *Signer
	}{
		{
			"success",
			args{
				sr25519test.CharlottePrivateKey,
			},
			&Signer{
				sr25519test.CharlottePrivateKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSigner(tt.args.privateKey); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewSigner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSigner_createSignature(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		privateKey crypto.PrivateKey
	}
	type args struct {
		signedData []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.MultiSignature
		wantErr bool
	}{
		{
			"sr25519-charlotte",
			fields{
				sr25519test.CharlottePrivateKey,
			},
			args{
				[]byte("signed-data"),
			},
			&types.MultiSignature{
				IsEd25519: false,
				AsEd25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				IsSr25519: true,
				AsSr25519: types.Signature{0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x2d, 0x64, 0x61, 0x74, 0x61, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				IsEcdsa:   false,
				AsEcdsa:   types.Bytes(nil),
			},
			false,
		},
		{
			"ed25519-charlotte",
			fields{
				ed25519test.CharlottePrivateKey,
			},
			args{
				[]byte("signed-data"),
			},
			&types.MultiSignature{
				IsEd25519: true,
				AsEd25519: types.Signature{0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x2d, 0x64, 0x61, 0x74, 0x61, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				IsSr25519: false,
				AsSr25519: types.Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				IsEcdsa:   false,
				AsEcdsa:   types.Bytes(nil),
			},
			false,
		},
		{
			"err-secp256k1-charlotte",
			fields{
				secp256k1test.CharlottePrivateKey,
			},
			args{
				[]byte("signed-data"),
			},
			nil,
			true,
		},
		{
			"err-invalid",
			fields{
				cryptotest.NewMockPrivateKey(mockCtrl),
			},
			args{
				[]byte("signed-data"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Signer{
				privateKey: tt.fields.privateKey,
			}
			got, err := e.createSignature(tt.args.signedData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signer.createSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Signer.createSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
