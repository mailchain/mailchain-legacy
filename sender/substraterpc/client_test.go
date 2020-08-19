package substraterpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mr-tron/base58/base58"

	gsrpc "github.com/mailchain/go-substrate-rpc-client"
	"github.com/mailchain/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetMetadata(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	type args struct {
		hash types.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Metadata
		wantErr bool
	}{
		{
			"success-latest",
			fields{
				func() *gsrpc.SubstrateAPI {
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV10String)))
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				types.Hash{},
			},
			types.ExamplaryMetadataV10,
			false,
		},
		{
			"success-specific",
			fields{
				func() *gsrpc.SubstrateAPI {
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV10String)))
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				types.NewHash([]byte("test")),
			},
			types.ExamplaryMetadataV10,
			false,
		},
		{
			"error-latest",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.WriteHeader(http.StatusBadRequest)
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				types.Hash{},
			},
			nil,
			true,
		},
		{
			"error-specific",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.WriteHeader(http.StatusBadRequest)
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				types.NewHash([]byte("test")),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.api)
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
				big.NewInt(SuggestedGas),
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
				big.NewInt(SuggestedGas),
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
			big.NewInt(SuggestedGas),
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

func TestSubstrateClient_GetBlockHash(t *testing.T) {
	assert := assert.New(t)
	var hash32 = []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2,
	}
	type args struct {
		blockNumber uint64
	}
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.Hash
		wantErr bool
	}{
		{
			"success-latest",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.Write([]byte("{\"result\":\"0x0102030405060708090001020304050607080900010203040506070809000102\"}"))
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				0,
			},
			types.NewHash(hash32),
			false,
		},
		{
			"success",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.Write([]byte("{\"result\":\"0x0102030405060708090001020304050607080900010203040506070809000102\"}"))
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				1,
			},
			types.NewHash(hash32),
			false,
		},
		{
			"error-latest",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.WriteHeader(http.StatusBadRequest)
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				0,
			},
			types.Hash{},
			true,
		},
		{
			"error",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.WriteHeader(http.StatusBadRequest)
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				1,
			},
			types.Hash{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.api)
			got, err := client.GetBlockHash(tt.args.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubstrateClient.GetBlockHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("SubstrateClient.GetBlockHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstrateClient_GetRuntimeVersion(t *testing.T) {
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	type args struct {
		blockHash types.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.RuntimeVersion
		wantErr error
	}{
		{
			"success-latest",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.Write([]byte("{\"jsonrpc\":\"2.0\",\"result\":{\"apis\":[[\"0xdf6acb689907609b\",2],[\"0x37e397fc7c91f5e4\",1],[\"0x40fe3ad401f8959a\",3],[\"0xd2bc9897eed08f15\",1],[\"0xf78b278be53f454c\",1],[\"0xed99c5acb25eedf5\",2],[\"0xbc9d89904f5b923f\",1],[\"0x687ad44ad37f03c2\",1],[\"0xdd718d5cc53262d4\",1],[\"0xab3c0572291feb8b\",1]],\"authoringVersion\":15,\"implName\":\"edgeware-node\",\"implVersion\":25,\"specName\":\"edgeware\",\"specVersion\":25},\"id\":3}"))
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				types.Hash{},
			},
			&types.RuntimeVersion{
				APIs: []types.RuntimeVersionAPI{
					{APIID: "0xdf6acb689907609b", Version: 0x2},
					{APIID: "0x37e397fc7c91f5e4", Version: 0x1},
					{APIID: "0x40fe3ad401f8959a", Version: 0x3},
					{APIID: "0xd2bc9897eed08f15", Version: 0x1},
					{APIID: "0xf78b278be53f454c", Version: 0x1},
					{APIID: "0xed99c5acb25eedf5", Version: 0x2},
					{APIID: "0xbc9d89904f5b923f", Version: 0x1},
					{APIID: "0x687ad44ad37f03c2", Version: 0x1},
					{APIID: "0xdd718d5cc53262d4", Version: 0x1},
					{APIID: "0xab3c0572291feb8b", Version: 0x1},
				},
				AuthoringVersion: 0xf,
				ImplName:         "edgeware-node",
				ImplVersion:      0x19,
				SpecName:         "edgeware",
				SpecVersion:      0x19,
			},
			nil,
		},
		{
			"success",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.Write([]byte("{\"jsonrpc\":\"2.0\",\"result\":{\"apis\":[[\"0xdf6acb689907609b\",2],[\"0x37e397fc7c91f5e4\",1],[\"0x40fe3ad401f8959a\",3],[\"0xd2bc9897eed08f15\",1],[\"0xf78b278be53f454c\",1],[\"0xed99c5acb25eedf5\",2],[\"0xbc9d89904f5b923f\",1],[\"0x687ad44ad37f03c2\",1],[\"0xdd718d5cc53262d4\",1],[\"0xab3c0572291feb8b\",1]],\"authoringVersion\":15,\"implName\":\"edgeware-node\",\"implVersion\":25,\"specName\":\"edgeware\",\"specVersion\":25},\"id\":3}"))
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				types.NewHash([]byte("test")),
			},
			&types.RuntimeVersion{
				APIs: []types.RuntimeVersionAPI{
					{APIID: "0xdf6acb689907609b", Version: 0x2},
					{APIID: "0x37e397fc7c91f5e4", Version: 0x1},
					{APIID: "0x40fe3ad401f8959a", Version: 0x3},
					{APIID: "0xd2bc9897eed08f15", Version: 0x1},
					{APIID: "0xf78b278be53f454c", Version: 0x1},
					{APIID: "0xed99c5acb25eedf5", Version: 0x2},
					{APIID: "0xbc9d89904f5b923f", Version: 0x1},
					{APIID: "0x687ad44ad37f03c2", Version: 0x1},
					{APIID: "0xdd718d5cc53262d4", Version: 0x1},
					{APIID: "0xab3c0572291feb8b", Version: 0x1},
				},
				AuthoringVersion: 0xf,
				ImplName:         "edgeware-node",
				ImplVersion:      0x19,
				SpecName:         "edgeware",
				SpecVersion:      0x19,
			},
			nil,
		},
		{
			"error-latest",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.WriteHeader(http.StatusBadRequest)
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				types.Hash{},
			},
			nil,
			errors.New("400 Bad Request "),
		},
		{
			"error",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.WriteHeader(http.StatusBadRequest)
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				types.NewHash([]byte("test")),
			},
			nil,
			errors.New("400 Bad Request "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.api)
			got, err := client.GetRuntimeVersion(tt.args.blockHash)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("SubstrateClient.GetRuntimeVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("SubstrateClient.GetRuntimeVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstrateClient_GetNonce(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		address  []byte
	}
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint32
		wantErr error
	}{
		{
			"error-pk",
			fields{
				func() *gsrpc.SubstrateAPI {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
					}))
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				context.Background(),
				"invalid",
				"edgeware",
				nil,
			},
			uint32(0),
			errors.New(`"invalid" unsupported protocol`),
		},
		{
			"error",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.WriteHeader(http.StatusBadRequest)
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				context.Background(),
				"substrate",
				"edgeware",
				func() []byte {
					addr, _ := base58.Decode("5CLmNK8f16nagFeF2h3iNeeChaxPiAsJu7piNYJgdPpmaRzP")
					return addr
				}(),
			},
			uint32(0),
			errors.New("400 Bad Request "),
		},
		{
			"success",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.Write([]byte("{\"result\":5}"))
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				context.Background(),
				"substrate",
				"edgeware",
				func() []byte {
					addr, _ := base58.Decode("5CLmNK8f16nagFeF2h3iNeeChaxPiAsJu7piNYJgdPpmaRzP")
					return addr
				}(),
			},
			uint32(5),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.api)
			got, err := client.GetNonce(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("SubstrateClient.GetBlockHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("SubstrateClient.GetBlockHash() = %v, want %v", got, tt.want)
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
		immortalEra bool
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
				true,
				func() types.RuntimeVersion {
					r := types.NewRuntimeVersion()
					return *r
				}(),
				123,
				1,
			},
			types.SignatureOptions{
				Era: types.ExtrinsicEra{
					IsImmortalEra: true,
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
			if got := s.CreateSignatureOptions(tt.args.blockHash, tt.args.genesisHash, tt.args.mortalEra, tt.args.immortalEra, tt.args.rv, tt.args.nonce, tt.args.tip); !assert.Equal(t, tt.want, got) {
				t.Errorf("SubstrateClient.CreateSignatureOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstrateClient_SubmitExtrinsic(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		xt *types.Extrinsic
	}
	type fields struct {
		api *gsrpc.SubstrateAPI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.Hash
		wantErr bool
	}{
		{
			"error",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.WriteHeader(http.StatusBadRequest)
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				&types.Extrinsic{Version: 0x84, Signature: types.ExtrinsicSignatureV4{Signer: types.Address{IsAccountID: true, AsAccountID: types.AccountID{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, IsAccountIndex: false, AsAccountIndex: 0x0}, Signature: types.MultiSignature{IsSr25519: true, AsSr25519: types.Signature{0xc0, 0x42, 0x19, 0x5f, 0x93, 0x25, 0xd, 0x3e, 0xda, 0xa2, 0xe4, 0xa4, 0x2d, 0xcf, 0x4e, 0x41, 0xc1, 0x6c, 0xa7, 0x1c, 0xfc, 0x3a, 0x2b, 0x23, 0x99, 0x8a, 0xd4, 0xec, 0x97, 0x4f, 0x8b, 0x1a, 0xcd, 0xcd, 0xad, 0x97, 0xd1, 0x4b, 0x6d, 0xf5, 0xcb, 0x89, 0x6, 0xff, 0x61, 0xc8, 0x92, 0x17, 0x96, 0x54, 0xa5, 0xec, 0xcc, 0xb, 0x66, 0x85, 0xf6, 0xc1, 0x7f, 0xed, 0x49, 0x21, 0x94, 0x0}}, Era: types.ExtrinsicEra{IsImmortalEra: true, IsMortalEra: false, AsMortalEra: types.MortalEra{First: 0x0, Second: 0x0}}, Nonce: 0x1, Tip: 0x0}, Method: types.Call{CallIndex: types.CallIndex{SectionIndex: 0x6, MethodIndex: 0x0}, Args: types.Args{0xff, 0x8e, 0xaf, 0x4, 0x15, 0x16, 0x87, 0x73, 0x63, 0x26, 0xc9, 0xfe, 0xa1, 0x7e, 0x25, 0xfc, 0x52, 0x87, 0x61, 0x36, 0x93, 0xc9, 0x12, 0x90, 0x9c, 0xb2, 0x26, 0xaa, 0x47, 0x94, 0xf2, 0x6a, 0x48, 0xe5, 0x6c}}},
			},
			types.Hash{},
			true,
		},
		{
			"success",
			fields{
				func() *gsrpc.SubstrateAPI {
					apiCall := true
					server := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if apiCall {
								apiCall = false
								w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
							} else {
								w.Write([]byte("{\"result\":\"0x9a8ef9794ded03b4d1ae45034351210e87f970f1f9500994bca82f9cd5a1166e\"}"))
							}
						}),
					)
					api, _ := gsrpc.NewSubstrateAPI(server.URL)
					return api
				}(),
			},
			args{
				&types.Extrinsic{Version: 0x84, Signature: types.ExtrinsicSignatureV4{Signer: types.Address{IsAccountID: true, AsAccountID: types.AccountID{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, IsAccountIndex: false, AsAccountIndex: 0x0}, Signature: types.MultiSignature{IsSr25519: true, AsSr25519: types.Signature{0xc0, 0x42, 0x19, 0x5f, 0x93, 0x25, 0xd, 0x3e, 0xda, 0xa2, 0xe4, 0xa4, 0x2d, 0xcf, 0x4e, 0x41, 0xc1, 0x6c, 0xa7, 0x1c, 0xfc, 0x3a, 0x2b, 0x23, 0x99, 0x8a, 0xd4, 0xec, 0x97, 0x4f, 0x8b, 0x1a, 0xcd, 0xcd, 0xad, 0x97, 0xd1, 0x4b, 0x6d, 0xf5, 0xcb, 0x89, 0x6, 0xff, 0x61, 0xc8, 0x92, 0x17, 0x96, 0x54, 0xa5, 0xec, 0xcc, 0xb, 0x66, 0x85, 0xf6, 0xc1, 0x7f, 0xed, 0x49, 0x21, 0x94, 0x0}}, Era: types.ExtrinsicEra{IsImmortalEra: true, IsMortalEra: false, AsMortalEra: types.MortalEra{First: 0x0, Second: 0x0}}, Nonce: 0x1, Tip: 0x0}, Method: types.Call{CallIndex: types.CallIndex{SectionIndex: 0x6, MethodIndex: 0x0}, Args: types.Args{0xff, 0x8e, 0xaf, 0x4, 0x15, 0x16, 0x87, 0x73, 0x63, 0x26, 0xc9, 0xfe, 0xa1, 0x7e, 0x25, 0xfc, 0x52, 0x87, 0x61, 0x36, 0x93, 0xc9, 0x12, 0x90, 0x9c, 0xb2, 0x26, 0xaa, 0x47, 0x94, 0xf2, 0x6a, 0x48, 0xe5, 0x6c}}},
			},
			types.Hash{0x9a, 0x8e, 0xf9, 0x79, 0x4d, 0xed, 0x3, 0xb4, 0xd1, 0xae, 0x45, 0x3, 0x43, 0x51, 0x21, 0xe, 0x87, 0xf9, 0x70, 0xf1, 0xf9, 0x50, 0x9, 0x94, 0xbc, 0xa8, 0x2f, 0x9c, 0xd5, 0xa1, 0x16, 0x6e},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.api)
			got, err := client.SubmitExtrinsic(tt.args.xt)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubstrateClient.SubmitExtrinsic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("SubstrateClient.SubmitExtrinsic() = %v, want %v", got, tt.want)
			}
		})
	}
}
