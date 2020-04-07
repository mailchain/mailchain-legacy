package substraterpc

import (
	"context"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/mailbox/signer/signertest"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/substraterpc/substraterpctest"
	"github.com/pkg/errors"
	"math/big"
	"testing"
)

func TestSubstrateRPC_Send(t *testing.T) {
	to := encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
	from := encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type signerOpts struct {
		arg1 string
	}
	type fields struct {
		client Client
	}
	type args struct {
		ctx     context.Context
		network string
		to      []byte
		from    []byte
		data    []byte
		signer  signer.Signer
		opts    sender.SendOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"success",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					metadata := &types.Metadata{
						MagicNumber:  0,
						Version:      0,
						IsMetadataV4: false,
						AsMetadataV4: types.MetadataV4{},
						IsMetadataV7: false,
						AsMetadataV7: types.MetadataV7{},
						IsMetadataV8: false,
						AsMetadataV8: types.MetadataV8{},
						IsMetadataV9: true,
						AsMetadataV9: types.MetadataV9{},
					}
					m.EXPECT().GetMetadata(types.Hash{}).Return(metadata, nil)
					addressTo := types.Address{
						IsAccountID:    true,
						AsAccountID:    types.AccountID{},
						IsAccountIndex: false,
						AsAccountIndex: 0,
					}
					m.EXPECT().GetAddress(to).Return(addressTo)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(SuggestedGas), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(SuggestedGas), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					rv := &types.RuntimeVersion{}
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(rv, nil)
					nonce := uint32(0)
					m.EXPECT().GetNonce(context.Background(), protocols.Substrate, substrate.EdgewareBerlin, from, metadata).Return(nonce, nil)
					m.EXPECT().CreateSignatureOptions(types.Hash{}, types.Hash{}, false, true, *rv, nonce, uint32(0))
					extrinsic := types.NewExtrinsic(types.Call{})
					m.EXPECT().SubmitExtrinsic(&extrinsic).Return(types.Hash{}, nil)
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      to,
				from:    from,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					signedExt := types.NewExtrinsic(types.Call{})
					m.EXPECT().Sign(substrate.SignerOptions{
						Extrinsic: types.NewExtrinsic(types.Call{}),
						SignatureOptions: types.SignatureOptions{
							BlockHash:   types.Hash{},
							Era:         types.ExtrinsicEra{IsMortalEra: false},
							GenesisHash: types.Hash{},
							Nonce:       types.UCompact(uint32(0)),
							SpecVersion: types.NewU32(0),
							Tip:         0,
						},
					}).Return(&signedExt, nil)
					return m
				}(),
				opts: signerOpts{"value1"},
			},
			false,
		},
		{
			"error-getMetadata",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					m.EXPECT().GetMetadata(types.Hash{}).Return(nil, errors.New("Get metadata error"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      nil,
				from:    nil,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					return signertest.NewMockSigner(mockCtrl)
				}(),
				opts: signerOpts{"value1"},
			},
			true,
		},
		{
			"error-estimate-gas",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					metadata := &types.Metadata{
						MagicNumber:  0,
						Version:      0,
						IsMetadataV4: false,
						AsMetadataV4: types.MetadataV4{},
						IsMetadataV7: false,
						AsMetadataV7: types.MetadataV7{},
						IsMetadataV8: false,
						AsMetadataV8: types.MetadataV8{},
						IsMetadataV9: true,
						AsMetadataV9: types.MetadataV9{},
					}
					m.EXPECT().GetMetadata(types.Hash{}).Return(metadata, nil)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(nil, errors.New("error suggest gas price"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      to,
				from:    nil,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					return signertest.NewMockSigner(mockCtrl)
				}(),
				opts: signerOpts{"value1"},
			},
			true,
		},
		{
			"error-call",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					metadata := &types.Metadata{
						MagicNumber:  0,
						Version:      0,
						IsMetadataV4: false,
						AsMetadataV4: types.MetadataV4{},
						IsMetadataV7: false,
						AsMetadataV7: types.MetadataV7{},
						IsMetadataV8: false,
						AsMetadataV8: types.MetadataV8{},
						IsMetadataV9: true,
						AsMetadataV9: types.MetadataV9{},
					}
					m.EXPECT().GetMetadata(types.Hash{}).Return(metadata, nil)
					addressTo := types.Address{
						IsAccountID:    true,
						AsAccountID:    types.AccountID{},
						IsAccountIndex: false,
						AsAccountIndex: 0,
					}
					m.EXPECT().GetAddress(to).Return(addressTo)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(SuggestedGas), nil)
					m.EXPECT().Call(metadata, addressTo, big.NewInt(SuggestedGas), []byte("transactionDataValue")).Return(types.Call{}, errors.New("Call error"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      to,
				from:    from,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					return signertest.NewMockSigner(mockCtrl)
				}(),
				opts: signerOpts{"value1"},
			},
			true,
		},
		{
			"error-get-block-hash",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					metadata := &types.Metadata{
						MagicNumber:  0,
						Version:      0,
						IsMetadataV4: false,
						AsMetadataV4: types.MetadataV4{},
						IsMetadataV7: false,
						AsMetadataV7: types.MetadataV7{},
						IsMetadataV8: false,
						AsMetadataV8: types.MetadataV8{},
						IsMetadataV9: true,
						AsMetadataV9: types.MetadataV9{},
					}
					m.EXPECT().GetMetadata(types.Hash{}).Return(metadata, nil)
					addressTo := types.Address{
						IsAccountID:    true,
						AsAccountID:    types.AccountID{},
						IsAccountIndex: false,
						AsAccountIndex: 0,
					}
					m.EXPECT().GetAddress(to).Return(addressTo)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(SuggestedGas), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(SuggestedGas), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, errors.New("error block hash"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      to,
				from:    from,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					return signertest.NewMockSigner(mockCtrl)
				}(),
				opts: signerOpts{"value1"},
			},
			true,
		},
		{
			"error-get-runtime-version",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					metadata := &types.Metadata{
						MagicNumber:  0,
						Version:      0,
						IsMetadataV4: false,
						AsMetadataV4: types.MetadataV4{},
						IsMetadataV7: false,
						AsMetadataV7: types.MetadataV7{},
						IsMetadataV8: false,
						AsMetadataV8: types.MetadataV8{},
						IsMetadataV9: true,
						AsMetadataV9: types.MetadataV9{},
					}
					m.EXPECT().GetMetadata(types.Hash{}).Return(metadata, nil)
					addressTo := types.Address{
						IsAccountID:    true,
						AsAccountID:    types.AccountID{},
						IsAccountIndex: false,
						AsAccountIndex: 0,
					}
					m.EXPECT().GetAddress(to).Return(addressTo)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(SuggestedGas), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(SuggestedGas), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(nil, errors.New("error runtime version"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      to,
				from:    from,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					return signertest.NewMockSigner(mockCtrl)
				}(),
				opts: signerOpts{"value1"},
			},
			true,
		},
		{
			"error-get-nonce",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					metadata := &types.Metadata{
						MagicNumber:  0,
						Version:      0,
						IsMetadataV4: false,
						AsMetadataV4: types.MetadataV4{},
						IsMetadataV7: false,
						AsMetadataV7: types.MetadataV7{},
						IsMetadataV8: false,
						AsMetadataV8: types.MetadataV8{},
						IsMetadataV9: true,
						AsMetadataV9: types.MetadataV9{},
					}
					m.EXPECT().GetMetadata(types.Hash{}).Return(metadata, nil)
					addressTo := types.Address{
						IsAccountID:    true,
						AsAccountID:    types.AccountID{},
						IsAccountIndex: false,
						AsAccountIndex: 0,
					}
					m.EXPECT().GetAddress(to).Return(addressTo)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(SuggestedGas), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(SuggestedGas), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(&types.RuntimeVersion{}, nil)
					m.EXPECT().GetNonce(context.Background(), protocols.Substrate, substrate.EdgewareBerlin, from, metadata).Return(uint32(0), errors.New("error get nonce"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      to,
				from:    from,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					return signertest.NewMockSigner(mockCtrl)
				}(),
				opts: signerOpts{"value1"},
			},
			true,
		},
		{
			"error-sign-tx",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					metadata := &types.Metadata{
						MagicNumber:  0,
						Version:      0,
						IsMetadataV4: false,
						AsMetadataV4: types.MetadataV4{},
						IsMetadataV7: false,
						AsMetadataV7: types.MetadataV7{},
						IsMetadataV8: false,
						AsMetadataV8: types.MetadataV8{},
						IsMetadataV9: true,
						AsMetadataV9: types.MetadataV9{},
					}
					m.EXPECT().GetMetadata(types.Hash{}).Return(metadata, nil)
					addressTo := types.Address{
						IsAccountID:    true,
						AsAccountID:    types.AccountID{},
						IsAccountIndex: false,
						AsAccountIndex: 0,
					}
					m.EXPECT().GetAddress(to).Return(addressTo)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(SuggestedGas), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(SuggestedGas), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					rv := &types.RuntimeVersion{}
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(rv, nil)
					nonce := uint32(0)
					m.EXPECT().GetNonce(context.Background(), protocols.Substrate, substrate.EdgewareBerlin, from, metadata).Return(nonce, nil)
					m.EXPECT().CreateSignatureOptions(types.Hash{}, types.Hash{}, false, true, *rv, nonce, uint32(0))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      to,
				from:    from,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					m.EXPECT().Sign(substrate.SignerOptions{
						Extrinsic: types.NewExtrinsic(types.Call{}),
						SignatureOptions: types.SignatureOptions{
							BlockHash:   types.Hash{},
							Era:         types.ExtrinsicEra{IsMortalEra: false},
							GenesisHash: types.Hash{},
							Nonce:       types.UCompact(uint32(0)),
							SpecVersion: types.NewU32(0),
							Tip:         0,
						},
					}).Return(nil, errors.New("error signing transaction"))
					return m
				}(),
				opts: signerOpts{"value1"},
			},
			true,
		},
		{
			"error-submit-tx",
			fields{
				func() Client {
					m := substraterpctest.NewMockClient(mockCtrl)
					metadata := &types.Metadata{
						MagicNumber:  0,
						Version:      0,
						IsMetadataV4: false,
						AsMetadataV4: types.MetadataV4{},
						IsMetadataV7: false,
						AsMetadataV7: types.MetadataV7{},
						IsMetadataV8: false,
						AsMetadataV8: types.MetadataV8{},
						IsMetadataV9: true,
						AsMetadataV9: types.MetadataV9{},
					}
					m.EXPECT().GetMetadata(types.Hash{}).Return(metadata, nil)
					addressTo := types.Address{
						IsAccountID:    true,
						AsAccountID:    types.AccountID{},
						IsAccountIndex: false,
						AsAccountIndex: 0,
					}
					m.EXPECT().GetAddress(to).Return(addressTo)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(SuggestedGas), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(SuggestedGas), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					rv := &types.RuntimeVersion{}
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(rv, nil)
					nonce := uint32(0)
					m.EXPECT().GetNonce(context.Background(), protocols.Substrate, substrate.EdgewareBerlin, from, metadata).Return(nonce, nil)
					m.EXPECT().CreateSignatureOptions(types.Hash{}, types.Hash{}, false, true, *rv, nonce, uint32(0))
					newExtrinsic := types.NewExtrinsic(types.Call{})
					m.EXPECT().SubmitExtrinsic(&newExtrinsic).Return(types.Hash{}, errors.New("error submitting transaction"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareBerlin,
				to:      to,
				from:    from,
				data:    []byte("transactionDataValue"),
				signer: func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					signedExt := types.NewExtrinsic(types.Call{})
					m.EXPECT().Sign(substrate.SignerOptions{
						Extrinsic: types.NewExtrinsic(types.Call{}),
						SignatureOptions: types.SignatureOptions{
							BlockHash:   types.Hash{},
							Era:         types.ExtrinsicEra{IsMortalEra: false},
							GenesisHash: types.Hash{},
							Nonce:       types.UCompact(uint32(0)),
							SpecVersion: types.NewU32(0),
							Tip:         0,
						},
					}).Return(&signedExt, nil)
					return m
				}(),
				opts: signerOpts{"value1"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubstrateRPC{
				client: tt.fields.client,
			}
			if err := s.Send(tt.args.ctx, tt.args.network, tt.args.to, tt.args.from, tt.args.data, tt.args.signer, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("SubstrateRPC.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
