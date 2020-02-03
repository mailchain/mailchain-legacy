package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/datastore/datastoretest"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type unknownPublicKey struct{}

func (pk unknownPublicKey) Bytes() []byte {
	return []byte("unknown public key")
}

func (pk unknownPublicKey) Kind() string {
	return "unknown"
}

func (pk unknownPublicKey) Verify(message, sig []byte) bool {
	return true
}

func Test_GetPublicKey(t *testing.T) {
	address := encodingtest.MustDecodeHexZeroX("0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6")
	var txHash = encodingtest.MustDecodeHexZeroX("0x98beb27135aa0a25650557005ad962919d6a278c4b3dde7f4f6a3a1e65aa746c")
	var blockHash = encodingtest.MustDecodeHexZeroX("0x373d339e45a701447367d7b9c7cef84aab79c2b2714271b908cda0ab3ad0849b")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		store datastore.PublicKeyStore
	}
	tests := []struct {
		name        string
		args        args
		queryParams map[string]string
		wantStatus  int
	}{
		{
			"422-invalid-request",
			args{
				nil,
			},
			map[string]string{},
			http.StatusUnprocessableEntity,
		},
		{
			"422-decode-by-protocol-unsupported-protocol",
			args{
				nil,
			},
			map[string]string{
				"address":  "0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6",
				"network":  "mainnet",
				"protocol": "unknown",
			},
			http.StatusUnprocessableEntity,
		},
		{
			"500-get-public-key-error",
			args{
				func() datastore.PublicKeyStore {
					store := datastoretest.NewMockPublicKeyStore(mockCtrl)
					store.EXPECT().GetPublicKey(gomock.Any(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).Return(nil, errors.New("error: GetPublicKey")).Times(1)
					return store
				}(),
			},
			map[string]string{
				"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusInternalServerError,
		},
		{
			"500-encryption-methods-unsupported-public-key-type",
			args{
				func() datastore.PublicKeyStore {
					store := datastoretest.NewMockPublicKeyStore(mockCtrl)
					res := &datastore.PublicKey{
						PublicKey: &unknownPublicKey{},
						BlockHash: blockHash,
						TxHash:    txHash,
					}
					store.EXPECT().GetPublicKey(gomock.Any(), "ethereum", "mainnet", address).Return(res, nil).Times(1)
					return store
				}(),
			},
			map[string]string{
				"address":  "0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusInternalServerError,
		},
		{
			"200-sofia-secp256k1",
			args{
				func() datastore.PublicKeyStore {
					store := datastoretest.NewMockPublicKeyStore(mockCtrl)
					res := &datastore.PublicKey{
						PublicKey: secp256k1test.SofiaPublicKey,
						BlockHash: blockHash,
						TxHash:    txHash,
					}
					store.EXPECT().GetPublicKey(gomock.Any(), "ethereum", "mainnet", address).Return(res, nil).Times(1)
					return store
				}(),
			},
			map[string]string{
				"address":  "0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusOK,
		},
		{
			"200-charlotte-secp256k1",
			args{
				func() datastore.PublicKeyStore {
					store := datastoretest.NewMockPublicKeyStore(mockCtrl)
					res := &datastore.PublicKey{
						PublicKey: secp256k1test.CharlottePublicKey,
						BlockHash: blockHash,
						TxHash:    txHash,
					}
					store.EXPECT().GetPublicKey(gomock.Any(), "ethereum", "mainnet", address).Return(res, nil).Times(1)
					return store
				}(),
			},
			map[string]string{
				"address":  "0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		testName := t.Name()
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetPublicKey(tt.args.store))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}

			golden, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s.json", testName, tt.name))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			assert.JSONEq(t, string(golden), rr.Body.String())
		})
	}
}
