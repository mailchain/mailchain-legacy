package etherscan

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewAPIClient(t *testing.T) {
	assert := assert.New(t)
	apiClient, _ := NewAPIClient("api-key")
	want := &APIClient{
		key: "api-key",
		networkConfigs: map[string]networkConfig{
			ethereum.Mainnet: {url: "https://api.etherscan.io/api"},
			ethereum.Ropsten: {url: "https://api-ropsten.etherscan.io/api"},
			ethereum.Kovan:   {url: "https://api-kovan.etherscan.io/api"},
			ethereum.Rinkeby: {url: "https://api-rinkeby.etherscan.io/api"},
			ethereum.Goerli:  {url: "https://api-goerli.etherscan.io/api"},
		},
	}
	if !assert.Equal(want, apiClient) {
		t.Errorf("NewAPIClient() = %v, want %v", apiClient, want)
	}
}

func TestGetTransactionByHash(t *testing.T) {
	assert := assert.New(t)
	networkStackError := errors.New("Get http://somethignnotvalid:1334")
	type args struct {
		server  *httptest.Server
		network string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		wantNil bool
		want    *types.Transaction
	}{
		{
			"success",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						txData := "{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"blockHash\":\"0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2\",\"blockNumber\":\"0x5daf3b\",\"from\":\"0xa7d9ddbe1f17865597fbd27ec712455208b6b76d\",\"gas\":\"0xc350\",\"gasPrice\":\"0x4a817c800\",\"hash\":\"0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b\",\"input\":\"0x68656c6c6f21\",\"nonce\":\"0x15\",\"to\":\"0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb\",\"transactionIndex\":\"0x41\",\"value\":\"0xf3dbb76162000\",\"v\":\"0x25\",\"r\":\"0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea\",\"s\":\"0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c\"}}"
						w.Write([]byte(txData))
					}),
				),
				"TestNetwork",
			},
			nil,
			false,
			func() *types.Transaction {
				return types.NewTransaction(
					uint64(21),
					common.HexToAddress("0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb"),
					big.NewInt(4290000000000000),
					uint64(50000),
					big.NewInt(int64(20000000000)),
					[]byte("hello!"))
			}(),
		},
		{
			"unsupported-network",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				),
				"UnsupportedNetwork",
			},
			errors.New("network not supported"),
			true,
			nil,
		},
		{
			"response-error",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
					)
					s.URL = "http://somethignnotvalid:1334"
					return s
				}(),
				"TestNetwork",
			},
			networkStackError,
			true,
			nil,
		},
		{
			"unmarshal-error",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				),
				"TestNetwork",
			},
			errors.New("unexpected end of JSON input"),
			true,
			nil,
		},
		{
			"error-body",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("{\"error\":{\"code\": 0, \"message\": \"error\"}}"))
					}),
				),
				"TestNetwork",
			},
			errors.New("error"),
			true,
			nil,
		},
		{
			"error-not-found",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("{\"blockHash\":\"0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2\",\"blockNumber\":\"0x5daf3b\",\"from\":\"0xa7d9ddbe1f17865597fbd27ec712455208b6b76d\",\"gas\":\"0xc350\",\"gasPrice\":\"0x4a817c800\",\"hash\":\"0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b\",\"input\":\"\",\"nonce\":\"0x15\",\"to\":\"0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb\",\"transactionIndex\":\"0x41\",\"value\":\"0xf3dbb76162000\",\"v\":\"0x25\",\"r\":\"0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea\",\"s\":\"0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c\"}"))
					}),
				),
				"TestNetwork",
			},
			errors.New("not found"),
			true,
			nil,
		},
		{
			"trx-unmarshall-error",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("{\"result\": \"plain\"}"))
					}),
				),
				"TestNetwork",
			},
			errors.New("json: cannot unmarshal string into Go value of type types.txdata"),
			true,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.args.server.Close()
			client := &APIClient{
				key:            "api-key",
				networkConfigs: map[string]networkConfig{"TestNetwork": {url: tt.args.server.URL}},
			}
			got, err := client.getTransactionByHash(tt.args.network, common.Hash{})
			if (err != nil) && tt.wantErr == networkStackError && !strings.HasPrefix(err.Error(), networkStackError.Error()) {
				t.Errorf("APIClient.getTransactionByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr != networkStackError && err.Error() != tt.wantErr.Error() {
				fmt.Print(err.Error())
				t.Errorf("APIClient.getTransactionByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("APIClient.getTransactionByHash() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
			if got != nil &&
				(!assert.Equal(tt.want.Nonce(), got.Nonce()) ||
					!assert.Equal(tt.want.To(), got.To()) ||
					!assert.Equal(tt.want.Value(), got.Value()) ||
					!assert.Equal(tt.want.Gas(), got.Gas()) ||
					!assert.Equal(tt.want.GasPrice(), got.GasPrice()) ||
					!assert.Equal(tt.want.Data(), got.Data())) {
				t.Errorf("APIClient.getTransactionByHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransactionsByAddress(t *testing.T) {
	assert := assert.New(t)
	networkStackError := errors.New("Get http://somethignnotvalid:1334")
	type args struct {
		server  *httptest.Server
		network string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		wantNil bool
		want    *txList
	}{
		{
			"success",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						txData := "{\"status\":\"1\",\"message\":\"OK\",\"result\":[{\"blockNumber\":\"65204\",\"timeStamp\":\"1439232889\",\"hash\":\"0x98beb27135aa0a25650557005ad962919d6a278c4b3dde7f4f6a3a1e65aa746c\",\"nonce\":\"0\",\"blockHash\":\"0x373d339e45a701447367d7b9c7cef84aab79c2b2714271b908cda0ab3ad0849b\",\"transactionIndex\":\"0\",\"from\":\"0x3fb1cd2cd96c6d5c0b5eb3322d807b34482481d4\",\"to\":\"0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae\",\"value\":\"0\",\"gas\":\"122261\",\"gasPrice\":\"50000000000\",\"isError\":\"0\",\"txreceipt_status\":\"\",\"input\":\"0xf00d4b5d000000000000000000000000036c8cecce8d8bbf0831d840d7f29c9e3ddefa63000000000000000000000000c5a96db085dda36ffbe390f455315d30d6d3dc52\",\"contractAddress\":\"\",\"cumulativeGasUsed\":\"122207\",\"gasUsed\":\"122207\",\"confirmations\":\"8881309\"},{\"blockNumber\":\"65342\",\"timeStamp\":\"1439235315\",\"hash\":\"0x621de9a006b56c425d21ee0e04ab25866fff4cf606dd5d03cf677c5eb2172161\",\"nonce\":\"1\",\"blockHash\":\"0x889d18b8791f43688d07e0b588e94de746a020d4337c61e5285cd97556a6416e\",\"transactionIndex\":\"0\",\"from\":\"0x3fb1cd2cd96c6d5c0b5eb3322d807b34482481d4\",\"to\":\"0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae\",\"value\":\"0\",\"gas\":\"122269\",\"gasPrice\":\"50000000000\",\"isError\":\"0\",\"txreceipt_status\":\"\",\"input\":\"0xf00d4b5d00000000000000000000000005096a47749d8bfab0a90c1bb7a95115dbe4cea60000000000000000000000005ed8cee6b63b1c6afce3ad7c92f4fd7e1b8fad9f\",\"contractAddress\":\"\",\"cumulativeGasUsed\":\"122207\",\"gasUsed\":\"122207\",\"confirmations\":\"8881171\"}]}"
						w.Write([]byte(txData))
					}),
				),
				"TestNetwork",
			},
			nil,
			false,
			&txList{
				Status:  "1",
				Message: "OK",
				Result: []txResult{
					{
						BlockNumber:       "65204",
						TimeStamp:         "1439232889",
						Hash:              "0x98beb27135aa0a25650557005ad962919d6a278c4b3dde7f4f6a3a1e65aa746c",
						Nonce:             "0",
						BlockHash:         "0x373d339e45a701447367d7b9c7cef84aab79c2b2714271b908cda0ab3ad0849b",
						TransactionIndex:  "0",
						From:              "0x3fb1cd2cd96c6d5c0b5eb3322d807b34482481d4",
						To:                "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
						Value:             "0",
						Gas:               "122261",
						GasPrice:          "50000000000",
						IsError:           "0",
						Input:             "0xf00d4b5d000000000000000000000000036c8cecce8d8bbf0831d840d7f29c9e3ddefa63000000000000000000000000c5a96db085dda36ffbe390f455315d30d6d3dc52",
						ContractAddress:   "",
						CumulativeGasUsed: "122207",
						GasUsed:           "122207",
						Confirmations:     "8881309",
					},
					{
						BlockNumber:       "65342",
						TimeStamp:         "1439235315",
						Hash:              "0x621de9a006b56c425d21ee0e04ab25866fff4cf606dd5d03cf677c5eb2172161",
						Nonce:             "1",
						BlockHash:         "0x889d18b8791f43688d07e0b588e94de746a020d4337c61e5285cd97556a6416e",
						TransactionIndex:  "0",
						From:              "0x3fb1cd2cd96c6d5c0b5eb3322d807b34482481d4",
						To:                "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
						Value:             "0",
						Gas:               "122269",
						GasPrice:          "50000000000",
						IsError:           "0",
						Input:             "0xf00d4b5d00000000000000000000000005096a47749d8bfab0a90c1bb7a95115dbe4cea60000000000000000000000005ed8cee6b63b1c6afce3ad7c92f4fd7e1b8fad9f",
						ContractAddress:   "",
						CumulativeGasUsed: "122207",
						GasUsed:           "122207",
						Confirmations:     "8881171",
					},
				},
			},
		},
		{
			"unsupported-network",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				),
				"UnsupportedNetwork",
			},
			errors.New("network not supported"),
			true,
			nil,
		},
		{
			"response-error",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
					)
					s.URL = "http://somethignnotvalid:1334"
					return s
				}(),
				"TestNetwork",
			},
			networkStackError,
			true,
			nil,
		},
		{
			"unmarshal-error",
			args{
				httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				),
				"TestNetwork",
			},
			errors.New("unexpected end of JSON input"),
			true,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.args.server.Close()
			client := &APIClient{
				key:            "api-key",
				networkConfigs: map[string]networkConfig{"TestNetwork": {url: tt.args.server.URL}},
			}
			got, err := client.getTransactionsByAddress(tt.args.network, []byte{})
			if (err != nil) && tt.wantErr == networkStackError && !strings.HasPrefix(err.Error(), networkStackError.Error()) {
				t.Errorf("APIClient.getTransactionsByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr != networkStackError && err.Error() != tt.wantErr.Error() {
				t.Errorf("APIClient.getTransactionsByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("APIClient.getTransactionsByAddress() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("APIClient.getTransactionsByAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
