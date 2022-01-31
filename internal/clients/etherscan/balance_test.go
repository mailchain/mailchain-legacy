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

package etherscan

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailchain/mailchain/stores"
	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	type args struct {
		ctx      context.Context
		protocol string
		network  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		wantNil bool
		want    []stores.Transaction
	}{
		{
			"err-network-not-supported",
			args{
				context.Background(),
				"ethereum",
				"InvalidNetwork",
			},
			errors.New("network not supported"),
			true,
			nil,
		},
		{
			"err-unmarshal",
			args{
				context.Background(),
				"ethereum",
				"TestNetwork",
			},
			errors.New("{invalid}: invalid character 'i' looking for beginning of object key string"),
			true,
			nil,
		},
		{
			"err-get",
			args{
				context.Background(),
				"ethereum",
				"TestNetwork",
			},
			nil,
			false,
			[]stores.Transaction{},
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
			got, err := client.Receive(tt.args.ctx, tt.args.protocol, tt.args.network, []byte{})
			if (err != nil) && !assert.Equal(t, tt.wantErr.Error(), err.Error()) {
				t.Errorf("APIClient.Receive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("APIClient.Receive() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("APIClient.Receive() = %v, want %v", got, tt.want)
			}
		})
	}
}
