package substrate

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/stretchr/testify/assert"
)

func TestSigner_Sign(t *testing.T) {
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
				SignerOptions{
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
			if err == nil {
				assert.IsType(t, &types.Extrinsic{}, gotSignedTransactionInterface)
				gotSignedTransaction := gotSignedTransactionInterface.(*types.Extrinsic)
				if !assert.Equal(t, tt.wantSigner, &gotSignedTransaction.Signature.Signer) {
					t.Errorf("Signer.Sign().Signature.Signer = %v, want %v", gotSignedTransaction.Signature.Signer, tt.wantSigner)
				}
				if !assert.Equal(t, tt.wantMethod, &gotSignedTransaction.Method) {
					t.Errorf("Signer.Sign().Method = %v, want %v", gotSignedTransaction.Method, tt.wantMethod)
				}
				switch tt.fields.privateKey.(type) {
				case *sr25519.PrivateKey:
					assert.NotEqual(t, [64]byte{}, gotSignedTransaction.Signature.Signature.AsSr25519[:])
				default:
					t.Error("unsupported key type")
				}
			}
		})
	}
}
