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
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(32000), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(32000), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					rv := &types.RuntimeVersion{}
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(rv, nil)
					nonce := uint32(0)
					m.EXPECT().GetNonce(context.Background(), protocols.Substrate, substrate.EdgewareTestnet, from, metadata).Return(nonce, nil)
					m.EXPECT().CreateSignatureOptions(types.Hash{}, types.Hash{}, false, *rv, nonce, uint32(0))
					m.EXPECT().SubmitExtrinsic(types.NewExtrinsic(types.Call{})).Return(types.Hash{}, nil)
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareTestnet,
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
					}).Return(types.NewExtrinsic(types.Call{}), nil)
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
				network: substrate.EdgewareTestnet,
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
				network: substrate.EdgewareTestnet,
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
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(32000), nil)
					m.EXPECT().Call(metadata, addressTo, big.NewInt(32000), []byte("transactionDataValue")).Return(types.Call{}, errors.New("Call error"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareTestnet,
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
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(32000), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(32000), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, errors.New("error block hash"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareTestnet,
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
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(32000), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(32000), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(nil, errors.New("error runtime version"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareTestnet,
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
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(32000), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(32000), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(&types.RuntimeVersion{}, nil)
					m.EXPECT().GetNonce(context.Background(), protocols.Substrate, substrate.EdgewareTestnet, from, metadata).Return(uint32(0), errors.New("error get nonce"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareTestnet,
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
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(32000), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(32000), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					rv := &types.RuntimeVersion{}
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(rv, nil)
					nonce := uint32(0)
					m.EXPECT().GetNonce(context.Background(), protocols.Substrate, substrate.EdgewareTestnet, from, metadata).Return(nonce, nil)
					m.EXPECT().CreateSignatureOptions(types.Hash{}, types.Hash{}, false, *rv, nonce, uint32(0))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareTestnet,
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
					m.EXPECT().SuggestGasPrice(context.Background()).Return(big.NewInt(32000), nil)
					call := types.Call{}
					m.EXPECT().Call(metadata, addressTo, big.NewInt(32000), []byte("transactionDataValue")).Return(call, nil)
					ext := types.NewExtrinsic(call)
					m.EXPECT().NewExtrinsic(call).Return(ext)
					m.EXPECT().GetBlockHash(uint64(0)).Return(types.Hash{}, nil)
					rv := &types.RuntimeVersion{}
					m.EXPECT().GetRuntimeVersion(types.Hash{}).Return(rv, nil)
					nonce := uint32(0)
					m.EXPECT().GetNonce(context.Background(), protocols.Substrate, substrate.EdgewareTestnet, from, metadata).Return(nonce, nil)
					m.EXPECT().CreateSignatureOptions(types.Hash{}, types.Hash{}, false, *rv, nonce, uint32(0))
					m.EXPECT().SubmitExtrinsic(types.NewExtrinsic(types.Call{})).Return(types.Hash{}, errors.New("error submitting transaction"))
					return m
				}(),
			},
			args{
				ctx:     context.Background(),
				network: substrate.EdgewareTestnet,
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
					}).Return(types.NewExtrinsic(types.Call{}), nil)
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

// // not real test, to be deleted
// func Test_SubstrateCallExample(t *testing.T) {
// 	t.Run("s", func(t *testing.T) {
// 		// This sample shows how to create a transaction to make a transfer from one an account to another.

// 		// Instantiate the API
// 		api, err := gsrpc.NewSubstrateAPI("ws://testnet3.edgewa.re:9944")
// 		if err != nil {
// 			panic(err)
// 		}

// 		meta, err := api.RPC.State.GetMetadataLatest()
// 		if err != nil {
// 			panic(err)
// 		}

// 		//from, err := types.NewAddressFromHexAccountID("0x9aff809089623c64ccd7ca79b8fffb114cef2c251b2b3e4f2d5144eee962fb6f")
// 		to, err := types.NewAddressFromHexAccountID("0xf49a6b6deb3cdb00d86e735ced8f48b649273fc6ff6a3bbba021a6b4a1c5d067") // account id
// 		if err != nil {
// 			panic(err)
// 		}

// 		c, err := types.NewCall(meta, "Contracts.call", to, types.UCompact(0), types.UCompact(32000), "0x6d61696c636861696e383162336636383539326431393338396439656432346664636338316331666630323835383962653535303436303532366631633961613436623864333739346337653032616565363563386631373733376361366637333564393565303965366131396636303838366638313239326535373835373133343562386531653466393238326531306433396637316238636639653731613231656336393939333637346634616261643231623831393531646565346665643565666465663334643131303264346333336538626662613330623461343730646162643434653938653262363439346136653862363963393336353864393631393639356633313561356266356262313865363265336266623237363463363335323631616366363730303862353761316262333838353164396132656635353730323861336166373839646537396234346662346130336137653637393037343030376531623237")
// 		if err != nil {
// 			panic(err)
// 		}

// 		// Create the extrinsic
// 		ext := types.NewExtrinsic(c)
// 		if err != nil {
// 			panic(err)
// 		}

// 		genesisHash, err := api.RPC.Chain.GetBlockHash(0)
// 		if err != nil {
// 			panic(err)
// 		}

// 		rv, err := api.RPC.State.GetRuntimeVersionLatest()
// 		if err != nil {
// 			panic(err)
// 		}
// 		pair, _ := signature.KeyringPairFromSecret("<private key or seed phrase here>")
// 		key, err := types.CreateStorageKey(meta, "System", "AccountNonce", pair.PublicKey, nil)
// 		if err != nil {
// 			panic(err)
// 		}

// 		var nonce uint32
// 		err = api.RPC.State.GetStorageLatest(key, &nonce)
// 		if err != nil {
// 			panic(err)
// 		}

// 		o := types.SignatureOptions{
// 			BlockHash:   genesisHash,
// 			Era:         types.ExtrinsicEra{IsMortalEra: false},
// 			GenesisHash: genesisHash,
// 			Nonce:       types.UCompact(nonce),
// 			SpecVersion: rv.SpecVersion,
// 			Tip:         0,
// 		}

// 		// Sign the transaction using Alice's default account
// 		err = ext.Sign(pair, o)
// 		if err != nil {
// 			panic(err)
// 		}

// 		// Send the extrinsic
// 		hash, err := api.RPC.Author.SubmitExtrinsic(ext)
// 		if err != nil {
// 			panic(err)
// 		}

// 		fmt.Printf("Transfer sent with hash %#x\n", hash)
// 	})
// }
