package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/keystoretest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetBalances(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		ks keystore.Store
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantBody   string
		wantStatus int
	}{
		{
			"422-missing-protocol",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=mainnet", nil),
			"{\"code\":422,\"message\":\"'protocol' must be specified exactly once\"}\n",
			http.StatusUnprocessableEntity,
		},
		{
			"422-missing-network",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?protocol=ethereum", nil),
			"{\"code\":422,\"message\":\"'network' must be specified exactly once\"}\n",
			http.StatusUnprocessableEntity,
		},
		{
			"500-keystore-error",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("ethereum", "mainnet").Return(
						nil,
						errors.Errorf("error getting address"),
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil),
			"{\"code\":500,\"message\":\"error getting address\"}\n",
			http.StatusInternalServerError,
		},
		{
			"200-empty-address",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("ethereum", "mainnet").Return(
						[][]byte{},
						nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil),
			"{\"addresses\":[]}\n",
			http.StatusOK,
		},
		{
			"200-substrate-address",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("substrate", "edgeware-beresheet").Return(
						[][]byte{
							encodingtest.MustDecodeHex("2a169a11721851f5dff3541dd5c4b0b478ac1cd092c9d5976e83daa0d03f26620c464b"),
							encodingtest.MustDecodeHex("2a84623e7252e41138af6904e1b02304c941625f39e5762589125dc1a2f2cf2e30e02a"),
						},
						nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=edgeware-beresheet&protocol=substrate", nil),
			"{\"addresses\":[\"5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761\"]}\n",
			http.StatusOK,
		},
		{
			"200-ethereum-single-address",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("ethereum", "mainnet").Return(
						[][]byte{encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")},
						nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil),
			"{\"addresses\":[\"5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761\"]}\n",
			http.StatusOK,
		},
		{
			"200-ethereum-multi-address",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("ethereum", "mainnet").Return(
						[][]byte{
							encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
							encodingtest.MustDecodeHex("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
						},
						nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil),
			"{\"addresses\":[\"5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761\",\"4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2\"]}\n",
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		testName := t.Name()
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetAddresses(tt.args.ks))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			goldenResponse, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/response-%s.json", testName, tt.name))
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			if !assert.JSONEq(t, string(goldenResponse), rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), goldenResponse)
			}
		})
	}
}
