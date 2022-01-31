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

package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/addressing/addressingtest"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/keystoretest"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/mailbox/mailboxtest"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/statemock"
	"github.com/mailchain/mailchain/stores/storestest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_GetMessages(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		receivers        map[string]mailbox.Receiver
		inbox            stores.State
		cache            stores.Cache
		ks               keystore.Store
		deriveKeyOptions multi.OptionsBuilders
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantStatus int
	}{
		{
			"200-zero-messages",
			args{
				inbox: func() stores.State {
					stateMock := statemock.NewMockState(mockCtrl)
					stateMock.EXPECT().GetTransactions("ethereum", "mainnet", addressingtest.EthereumBob, int32(0), int32(15)).Times(1)
					return stateMock
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress(addressingtest.EthereumBob, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", fmt.Sprintf("/?address=%s&network=mainnet&protocol=ethereum", encoding.EncodeHexZeroX(addressingtest.EthereumBob)), nil),
			http.StatusOK,
		},
		{
			"200-message",
			args{
				cache: func() stores.Cache {
					m := storestest.NewMockCache(mockCtrl)
					m.EXPECT().GetMessage("test://TestReadMessage/success-2204f3d89e5a").Return([]byte{}, errors.Errorf("message not found"))
					m.EXPECT().SetMessage("test://TestReadMessage/success-2204f3d89e5a", gomock.Any()).Return(nil)
					return m
				}(),
				inbox: func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					inbox.EXPECT().GetReadStatus(mail.ID{71, 236, 160, 17, 227, 43, 82, 199, 16, 5, 173, 138, 143, 117, 225, 180, 76, 146, 201, 159, 209, 46, 67, 188, 207, 229, 113, 227, 194, 209, 61, 46, 154, 130, 106, 85, 15, 95, 246, 59, 36, 122, 244, 113}).Return(false, nil).Times(1)
					inbox.EXPECT().GetTransactions("ethereum", "mainnet", addressingtest.EthereumBob, int32(0), int32(15)).Return([]stores.Transaction{
						{
							BlockNumber:  100,
							Hash:         []byte{0x01, 0x02, 0x03},
							EnvelopeData: encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
						},
						{ // bad envelope
							BlockNumber:  102,
							Hash:         []byte{0x01, 0x02, 0x08},
							EnvelopeData: encodingtest.MustDecodeHex("500801127365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
						},
					}, nil)

					return inbox
				}(),
				ks: func() keystore.Store {
					decrypted, _ := ioutil.ReadFile("./testdata/simple.golden.eml")
					decrypter := ciphertest.NewMockDecrypter(mockCtrl)
					gomock.InOrder(
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent([]byte{0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65})).Return(decrypted, nil),
					)

					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress(addressingtest.EthereumBob, "ethereum", "mainnet").Return(true).Times(1)
					store.EXPECT().GetDecrypter(addressingtest.EthereumBob, "ethereum", "mainnet", byte(0x73), multi.OptionsBuilders{}).Return(decrypter, nil).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", fmt.Sprintf("/?address=%s&network=mainnet&protocol=ethereum", encoding.EncodeHexZeroX(addressingtest.EthereumBob)), nil),
			http.StatusOK,
		},
		{
			"200-message-rekeyed",
			args{
				cache: func() stores.Cache {
					m := storestest.NewMockCache(mockCtrl)
					m.EXPECT().GetMessage("test://TestReadMessage/success-2204f3d89e5a").Return([]byte{}, errors.Errorf("message not found"))
					m.EXPECT().SetMessage("test://TestReadMessage/success-2204f3d89e5a", gomock.Any()).Return(nil)
					return m
				}(),
				inbox: func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					inbox.EXPECT().GetReadStatus(mail.ID{71, 236, 160, 17, 227, 43, 82, 199, 16, 5, 173, 138, 143, 117, 225, 180, 76, 146, 201, 159, 209, 46, 67, 188, 207, 229, 113, 227, 194, 209, 61, 46, 154, 130, 106, 85, 15, 95, 246, 59, 36, 122, 244, 113}).Return(false, nil).Times(1)
					inbox.EXPECT().GetTransactions("ethereum", "mainnet", addressingtest.EthereumBob, int32(0), int32(15)).Return([]stores.Transaction{
						{
							BlockNumber:  100,
							Hash:         []byte{0x01, 0x02, 0x03},
							EnvelopeData: encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
							RekeyAddress: addressingtest.EthereumBob,
						},
					}, nil)
					return inbox
				}(),
				ks: func() keystore.Store {
					decrypted, _ := ioutil.ReadFile("./testdata/simple.golden.eml")
					decrypter := ciphertest.NewMockDecrypter(mockCtrl)
					gomock.InOrder(
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent([]byte{0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65})).Return(decrypted, nil),
					)

					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress(addressingtest.EthereumBob, "ethereum", "mainnet").Return(true).Times(1)
					store.EXPECT().GetDecrypter(addressingtest.EthereumBob, "ethereum", "mainnet", byte(0x73), multi.OptionsBuilders{}).Return(decrypter, nil).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", fmt.Sprintf("/?address=%s&network=mainnet&protocol=ethereum", encoding.EncodeHexZeroX(addressingtest.EthereumBob)), nil),
			http.StatusOK,
		},
		{
			"200-message-with-fetch",
			args{
				cache: func() stores.Cache {
					m := storestest.NewMockCache(mockCtrl)
					m.EXPECT().GetMessage("test://TestReadMessage/success-2204f3d89e5a").Return([]byte{}, errors.Errorf("message not found"))
					m.EXPECT().SetMessage("test://TestReadMessage/success-2204f3d89e5a", gomock.Any()).Return(nil)
					return m
				}(),
				inbox: func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					inbox.EXPECT().GetReadStatus(mail.ID{71, 236, 160, 17, 227, 43, 82, 199, 16, 5, 173, 138, 143, 117, 225, 180, 76, 146, 201, 159, 209, 46, 67, 188, 207, 229, 113, 227, 194, 209, 61, 46, 154, 130, 106, 85, 15, 95, 246, 59, 36, 122, 244, 113}).Return(false, nil).Times(1)
					inbox.EXPECT().GetTransactions("ethereum", "mainnet", addressingtest.EthereumBob, int32(0), int32(15)).Return([]stores.Transaction{
						{
							BlockNumber:  100,
							Hash:         []byte{0x01, 0x02, 0x03},
							EnvelopeData: encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
						},
					}, nil)
					inbox.EXPECT().PutTransaction("ethereum", "mainnet", addressingtest.EthereumBob, gomock.Any()).Return(nil).Times(1)

					return inbox
				}(),
				ks: func() keystore.Store {
					decrypted, _ := ioutil.ReadFile("./testdata/simple.golden.eml")
					decrypter := ciphertest.NewMockDecrypter(mockCtrl)
					gomock.InOrder(
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent([]byte{0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65})).Return(decrypted, nil),
					)

					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress(addressingtest.EthereumBob, "ethereum", "mainnet").Return(true).Times(1)
					store.EXPECT().GetDecrypter(addressingtest.EthereumBob, "ethereum", "mainnet", byte(0x73), multi.OptionsBuilders{}).Return(decrypter, nil).Times(1)
					return store
				}(),
				receivers: map[string]mailbox.Receiver{
					"ethereum/mainnet": func() mailbox.Receiver {
						receiver := mailboxtest.NewMockReceiver(mockCtrl)
						receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", addressingtest.EthereumBob).
							Return([]stores.Transaction{
								{
									EnvelopeData: encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
									BlockNumber:  100,
									Hash:         []byte("YS1oYXNo"),
								},
							}, nil).Times(1)
						return receiver
					}(),
				},
			},
			httptest.NewRequest("GET", fmt.Sprintf("/?address=%s&network=mainnet&protocol=ethereum&fetch=true", encoding.EncodeHexZeroX(addressingtest.EthereumBob)), nil),
			http.StatusOK,
		},
		{
			"500-get-transactions",
			args{
				inbox: func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					inbox.EXPECT().GetTransactions("ethereum", "mainnet", addressingtest.EthereumBob, int32(0), int32(15)).Return(nil, errors.New("internal error"))
					return inbox
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress(addressingtest.EthereumBob, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", fmt.Sprintf("/?address=%s&network=mainnet&protocol=ethereum", encoding.EncodeHexZeroX(addressingtest.EthereumBob)), nil),
			http.StatusInternalServerError,
		},
		{
			"406-no-private-key-found",
			args{
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress(addressingtest.EthereumBob, "ethereum", "mainnet").Return(false).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", fmt.Sprintf("/?address=%s&network=mainnet&protocol=ethereum", encoding.EncodeHexZeroX(addressingtest.EthereumBob)), nil),
			http.StatusNotAcceptable,
		},
		{
			"422-empty-address",
			args{},
			httptest.NewRequest("GET", "/?address=&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
	}

	testName := t.Name()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetMessages(tt.args.receivers, tt.args.inbox, tt.args.cache, tt.args.ks, tt.args.deriveKeyOptions))

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

func Test_parseGetMessagesRequest(t *testing.T) {
	type args struct {
		queryParams map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *GetMessagesRequest
		wantErr bool
	}{
		{
			"success-fetch",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "ethereum",
					"fetch":    "true",
				},
			},
			&GetMessagesRequest{
				Address:      "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				Network:      "mainnet",
				Protocol:     "ethereum",
				Fetch:        true,
				Limit:        15,
				Skip:         0,
				addressBytes: []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			},
			false,
		},
		{
			"success-no-fetch",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "ethereum",
				},
			},
			&GetMessagesRequest{
				Address:      "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				Network:      "mainnet",
				Protocol:     "ethereum",
				Limit:        15,
				Skip:         0,
				addressBytes: []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			},
			false,
		},
		{
			"err-skip-value",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "ethereum",
					"skip":     "invalid",
				},
			},
			nil,
			true,
		},
		{
			"err-limit-value",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "ethereum",
					"limit":    "invalid",
				},
			},
			nil,
			true,
		},
		{
			"err-invalid-protocol",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "invalid",
				},
			},
			nil,
			true,
		},
		{
			"err-empty-address",
			args{
				map[string]string{
					"address":  "",
					"network":  "mainnet",
					"protocol": "ethereum",
				},
			},
			nil,
			true,
		},
		{
			"err-empty-protocol",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "",
				},
			},
			nil,
			true,
		},
		{
			"err-empty-network",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "",
					"protocol": "ethereum",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for k, v := range tt.args.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			got, err := parseGetMessagesRequest(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGetMessagesRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseGetMessagesRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fetchMessages(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		ctx       context.Context
		protocol  string
		network   string
		address   []byte
		receivers map[string]mailbox.Receiver
		inbox     stores.State
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"single-transaction",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				addressingtest.EthereumBob,
				func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", addressingtest.EthereumBob).
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
				func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					inbox.EXPECT().PutTransaction("ethereum", "mainnet", addressingtest.EthereumBob, stores.Transaction{EnvelopeData: encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"), BlockNumber: 100, Hash: []byte("YS1oYXNo")}).Return(nil).Times(1)
					return inbox
				}(),
			},
			false,
		},
		{
			"zero-transactions",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				addressingtest.EthereumBob,
				func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", addressingtest.EthereumBob).Return([]stores.Transaction{}, nil).Times(1)
							return receiver
						}(),
					}
				}(),
				func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					return inbox
				}(),
			},
			false,
		},
		{
			"err-put-transaction",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				addressingtest.EthereumBob,
				func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", addressingtest.EthereumBob).
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
				func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					inbox.EXPECT().PutTransaction("ethereum", "mainnet", addressingtest.EthereumBob, stores.Transaction{EnvelopeData: encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"), BlockNumber: 100, Hash: []byte("YS1oYXNo")}).Return(errors.New("error")).Times(1)
					return inbox
				}(),
			},
			true,
		},
		{
			"err-receiver-network-not-supported",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				addressingtest.EthereumBob,
				func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", addressingtest.EthereumBob).Return([]stores.Transaction{}, errors.New("network not supported")).Times(1)
							return receiver
						}(),
					}
				}(),
				func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					return inbox
				}(),
			},
			true,
		},
		{
			"err-receiver",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				addressingtest.EthereumBob,
				func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", addressingtest.EthereumBob).Return([]stores.Transaction{}, errors.New("error")).Times(1)
							return receiver
						}(),
					}
				}(),
				func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					return inbox
				}(),
			},
			true,
		},
		{
			"err-nil-receiver",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				addressingtest.EthereumBob,
				func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": nil,
					}
				}(),
				func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					return inbox
				}(),
			},
			true,
		},
		{
			"err-no-receiver",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				addressingtest.EthereumBob,
				func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{}
				}(),
				func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					return inbox
				}(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fetchMessages(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address, tt.args.receivers, tt.args.inbox); (err != nil) != tt.wantErr {
				t.Errorf("fetchMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
