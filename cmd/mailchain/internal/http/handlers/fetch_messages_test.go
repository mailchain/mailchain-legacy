package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/addressing/addressingtest"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/keystoretest"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/mailbox/mailboxtest"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/statemock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_FetchMessages(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		inbox     stores.State
		receivers map[string]mailbox.Receiver
		ks        keystore.Store
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantStatus int
	}{
		{
			"422-empty-address",
			args{},
			httptest.NewRequest("GET", "/?address=&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"422-receiver-not-supported",
			args{},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"422-receiver-no-configured",
			args{
				receivers: map[string]mailbox.Receiver{
					"ethereum/mainnet": nil,
				},
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"406-no-private-key-found",
			args{
				receivers: map[string]mailbox.Receiver{
					"ethereum/mainnet": etherscan.APIClient{},
				},
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(false).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil).WithContext(context.Background()),
			http.StatusNotAcceptable,
		},
		{
			"406-receiver-network-error",
			args{
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).
								Return(nil, errors.New("network not supported")).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusNotAcceptable,
		},
		{
			"500-receiver-internal-error",
			args{
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).
								Return(nil, errors.New("internal error")).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusInternalServerError,
		},
		{
			"200-message",
			args{
				inbox: func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					inbox.EXPECT().PutTransaction("ethereum", "mainnet", addressingtest.EthereumCharlotte, stores.Transaction{EnvelopeData: encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"), BlockNumber: 100, Hash: []byte("YS1oYXNo")}).Return(nil).Times(1)
					return inbox
				}(),
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", addressingtest.EthereumCharlotte).
								Return([]stores.Transaction{
									{
										EnvelopeData: encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
										BlockNumber:  100,
										Hash:         []byte("YS1oYXNo"),
									},
								}, nil).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress(addressingtest.EthereumCharlotte, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", fmt.Sprintf("/?address=%s&network=mainnet&protocol=ethereum", encoding.EncodeHexZeroX(addressingtest.EthereumCharlotte)), nil),
			http.StatusOK,
		},
		{
			"200-zero-messages",
			args{
				inbox: func() stores.State {
					return statemock.NewMockState(mockCtrl)
				}(),
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).
								Return([]stores.Transaction{}, nil).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(FetchMessages(tt.args.inbox, tt.args.receivers, tt.args.ks))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
		})
	}
}
