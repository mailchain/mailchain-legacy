package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/keystoretest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetAddresses(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		ks keystore.Store
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantStatus int
	}{
		{
			"200-missing-protocol",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("", "").Return(
						map[string]map[string][][]uint8{
							"algorand": {
								"betanet": [][]uint8{}, "mainnet": [][]uint8{}, "testnet": [][]uint8{},
							},
							"ethereum": {
								"goerli":  [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
								"kovan":   [][]uint8{},
								"mainnet": [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
								"rinkeby": [][]uint8{},
								"ropsten": [][]uint8{}},
							"substrate": {
								"edgeware-beresheet": [][]uint8{
									{0x2a, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x6d, 0x48},
									{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x63, 0x83},
								},
								"edgeware-local": [][]uint8{},
								"edgeware-mainnet": [][]uint8{
									{0x7, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x9b, 0x76},
									{0x7, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0xda, 0xb},
								},
							},
						},
						nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/", nil),
			http.StatusOK,
		},
		{
			"200-missing-network",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("ethereum", "").Return(
						map[string]map[string][][]uint8{
							"ethereum": {
								"goerli":  [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
								"kovan":   [][]uint8{},
								"mainnet": [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
								"rinkeby": [][]uint8{},
								"ropsten": [][]uint8{}},
						},
						nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?protocol=ethereum", nil),
			http.StatusOK,
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
			http.StatusInternalServerError,
		},
		{
			"200-empty-address",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("ethereum", "mainnet").Return(
						map[string]map[string][][]uint8{"ethereum": {"mainnet": [][]uint8{}}},
						nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil),
			http.StatusOK,
		},
		{
			"200-substrate-edgeware-beresheet-address",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("substrate", "edgeware-beresheet").Return(
						map[string]map[string][][]uint8{
							"substrate": {
								"edgeware-beresheet": [][]uint8{
									{0x2a, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x6d, 0x48},
									{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x63, 0x83},
								},
							},
						},
						nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=edgeware-beresheet&protocol=substrate", nil),
			http.StatusOK,
		},
		{
			"200-ethereum-single-address",
			args{
				func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().GetAddresses("ethereum", "mainnet").Return(
						map[string]map[string][][]uint8{
							"ethereum": {
								"mainnet": [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
							},
						}, nil,
					).Times(1)

					return store
				}(),
			},
			httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil),
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
