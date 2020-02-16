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

package blockscout

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIClient(t *testing.T) {
	apiClient, _ := NewAPIClient("api-key")
	want := &APIClient{
		key: "api-key",
		networkConfigs: map[string]networkConfig{
			ethereum.Mainnet: {url: "https://blockscout.com/eth/mainnet/api"},
		},
	}
	if !assert.Equal(t, want, apiClient) {
		t.Errorf("NewAPIClient() = %v, want %v", apiClient, want)
	}
}

func TestGetTransactionsByAddress(t *testing.T) {
	type args struct {
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
			"err-unsupported-network",
			args{
				"UnsupportedNetwork",
			},
			errors.New("network not supported"),
			true,
			nil,
		},
		{
			"err-get",
			args{
				"TestNetwork",
			},
			errors.New("Invalid address format"),
			true,
			nil,
		},
		{
			"err-unmarshal",
			args{
				"TestNetwork",
			},
			errors.New("unexpected end of JSON input"),
			true,
			nil,
		},
	}
	for _, tt := range tests {
		testName := t.Name()
		t.Run(tt.name, func(t *testing.T) {
			golden, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s.json", testName, tt.name))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(golden))
				}),
			)
			defer server.Close()
			client := &APIClient{
				key:            "api-key",
				networkConfigs: map[string]networkConfig{"TestNetwork": {url: server.URL}},
			}
			got, err := client.getTransactionsByAddress(tt.args.network, []byte{})
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("APIClient.getTransactionsByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("APIClient.getTransactionsByAddress() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("APIClient.getTransactionsByAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
