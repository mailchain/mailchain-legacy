package substraterpc

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetMetadata(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		server *httptest.Server
		hash   types.Hash
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Metadata
		wantErr bool
	}{
		{
			"success-latest",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV4String)))
						}),
					)
					return s
				}(),
				types.Hash{},
			},
			types.ExamplaryMetadataV4,
			false,
		},
		{
			"success-specific",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV4String)))
						}),
					)
					return s
				}(),
				types.NewHash([]byte("test")),
			},
			types.ExamplaryMetadataV4,
			false,
		},
		{
			"error-latest",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusBadRequest)
						}),
					)
					return s
				}(),
				types.Hash{},
			},
			nil,
			true,
		},
		{
			"error-specific",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusBadRequest)
						}),
					)
					return s
				}(),
				types.NewHash([]byte("test")),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, _ := gsrpc.NewSubstrateAPI(tt.args.server.URL)
			client := NewClient(api)
			got, err := client.GetMetadata(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("Client.GetMetadata() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstrateClient_Call(t *testing.T) {
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	type args struct {
		metadata *types.Metadata
		to       types.Address
		gas      *big.Int
		data     []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.Call
		wantErr bool
	}{
		{
			"success",
			fields{},
			args{
				types.ExamplaryMetadataV8,
				func() types.Address {
					a := types.NewAddressFromAccountID(sr25519test.CharlottePublicKey.Bytes())
					return a
				}(),
				big.NewInt(32000),
				[]byte("message"),
			},
			types.Call{
				CallIndex: types.CallIndex{
					SectionIndex: 0x11,
					MethodIndex:  0x2,
				},
				Args: types.Args{
					0xff, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0x0, 0x2, 0xf4, 0x1, 0x0, 0x1c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
				},
			},
			false,
		},
		{
			"err-no-contracts",
			fields{},
			args{
				types.ExamplaryMetadataV4,
				types.NewAddressFromAccountID(sr25519test.CharlottePublicKey.Bytes()),
				big.NewInt(32000),
				[]byte("message"),
			},
			types.Call{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubstrateClient{
				api: tt.fields.api,
			}
			got, err := s.Call(tt.args.metadata, tt.args.to, tt.args.gas, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubstrateClient.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("SubstrateClient.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstrateClient_NewExtrinsic(t *testing.T) {
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	type args struct {
		call types.Call
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   types.Extrinsic
	}{
		{
			"success",
			fields{},
			args{
				types.Call{
					CallIndex: types.CallIndex{
						SectionIndex: 0x11,
						MethodIndex:  0x2,
					},
					Args: types.Args{
						0xff, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0x0, 0x2, 0xf4, 0x1, 0x0, 0x1c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
					},
				},
			},
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubstrateClient{
				api: tt.fields.api,
			}
			if got := s.NewExtrinsic(tt.args.call); !assert.Equal(t, tt.want, got) {
				t.Errorf("SubstrateClient.NewExtrinsic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstrateClient_GetAddress(t *testing.T) {
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	type args struct {
		accountID []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   types.Address
	}{
		{
			"success",
			fields{},
			args{
				sr25519test.CharlottePrivateKey.Bytes(),
			},
			types.Address{
				IsAccountID:    true,
				AsAccountID:    types.AccountID{0x23, 0xb0, 0x63, 0xa5, 0x81, 0xfd, 0x8e, 0x5e, 0x84, 0x7c, 0x4e, 0x2b, 0x9c, 0x49, 0x42, 0x47, 0x29, 0x87, 0x91, 0x53, 0xf, 0x52, 0x93, 0xbe, 0x36, 0x9e, 0x8b, 0xf2, 0x3a, 0x45, 0xd2, 0xbd},
				IsAccountIndex: false,
				AsAccountIndex: 0x0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubstrateClient{
				api: tt.fields.api,
			}
			if got := s.GetAddress(tt.args.accountID); !assert.Equal(t, tt.want, got) {
				t.Errorf("SubstrateClient.GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstrateClient_SuggestGasPrice(t *testing.T) {
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Int
		wantErr bool
	}{
		{
			"success",
			fields{},
			args{
				context.Background(),
			},
			big.NewInt(32000),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubstrateClient{
				api: tt.fields.api,
			}
			got, err := s.SuggestGasPrice(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubstrateClient.SuggestGasPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("SubstrateClient.SuggestGasPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstrateClient_CreateSignatureOptions(t *testing.T) {
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	type args struct {
		blockHash   types.Hash
		genesisHash types.Hash
		mortalEra   bool
		rv          types.RuntimeVersion
		nonce       uint32
		tip         uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   types.SignatureOptions
	}{
		{
			"success",
			fields{},
			args{
				types.Hash{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x30, 0x31, 0x32},
				types.Hash{0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0x30, 0xd1, 0xd2},
				false,
				func() types.RuntimeVersion {
					r := types.NewRuntimeVersion()
					return *r
				}(),
				123,
				1,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubstrateClient{
				api: tt.fields.api,
			}
			if got := s.CreateSignatureOptions(tt.args.blockHash, tt.args.genesisHash, tt.args.mortalEra, tt.args.rv, tt.args.nonce, tt.args.tip); !assert.Equal(t, tt.want, got) {
				t.Errorf("SubstrateClient.CreateSignatureOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
