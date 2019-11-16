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

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func TestGetProtocols(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		base *settings.Root
	}
	tests := []struct {
		name       string
		args       args
		wantBody   string
		wantStatus int
	}{
		{
			"disabled-ethereum",
			args{
				func() *settings.Root {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("protocols.ethereum.disabled").Return(true)
					m.EXPECT().GetBool("protocols.ethereum.disabled").Return(true)
					m.EXPECT().IsSet("protocols.substrate.disabled").Return(true)
					m.EXPECT().GetBool("protocols.substrate.disabled").Return(true)
					m.EXPECT().IsSet(gomock.Any()).Return(false).AnyTimes()
					return settings.FromStore(m)
				}(),
			},
			"{\"protocols\":[]}\n",
			http.StatusOK,
		},
		{
			"disabled-goreli",
			args{
				func() *settings.Root {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("protocols.ethereum.networks.goerli.disabled").Return(true)
					m.EXPECT().GetBool("protocols.ethereum.networks.goerli.disabled").Return(true)
					m.EXPECT().IsSet("protocols.substrate.disabled").Return(true)
					m.EXPECT().GetBool("protocols.substrate.disabled").Return(true)
					m.EXPECT().IsSet(gomock.Any()).Return(false).AnyTimes()
					return settings.FromStore(m)
				}(),
			},
			"{\"protocols\":[{\"name\":\"ethereum\",\"networks\":[\"kovan\",\"mainnet\",\"rinkeby\",\"ropsten\"]}]}\n",
			http.StatusOK,
		},
		{
			"default-ethereum",
			args{
				func() *settings.Root {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("protocols.substrate.disabled").Return(true)
					m.EXPECT().GetBool("protocols.substrate.disabled").Return(true)
					m.EXPECT().IsSet(gomock.Any()).Return(false).AnyTimes()
					return settings.FromStore(m)
				}(),
			},
			"{\"protocols\":[{\"name\":\"ethereum\",\"networks\":[\"goerli\",\"kovan\",\"mainnet\",\"rinkeby\",\"ropsten\"]}]}\n",
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetProtocols(tt.args.base))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			if !assert.Equal(tt.wantBody, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantBody)
			}
		})
	}
}
