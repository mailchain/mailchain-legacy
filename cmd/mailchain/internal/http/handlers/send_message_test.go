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
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/address/addresstest"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/stores"
	"github.com/stretchr/testify/assert"
)

func Test_checkForEmpties(t *testing.T) {
	type args struct {
		msg PostMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{
				PostMessage{
					Headers:   &PostHeaders{},
					Subject:   "subject-value",
					Body:      "body-value",
					PublicKey: "public-key-value",
				},
			},
			false,
		},
		{
			"empty-headers",
			args{
				PostMessage{
					Subject:   "subject-value",
					Body:      "body-value",
					PublicKey: "public-key-value",
				},
			},
			true,
		},
		{
			"empty-subject",
			args{
				PostMessage{
					Headers:   &PostHeaders{},
					Body:      "body-value",
					PublicKey: "public-key-value",
				},
			},
			true,
		},
		{
			"empty-body",
			args{
				PostMessage{
					Headers:   &PostHeaders{},
					Subject:   "subject-value",
					PublicKey: "public-key-value",
				},
			},
			true,
		},
		{
			"empty-public-key",
			args{
				PostMessage{
					Headers: &PostHeaders{},
					Subject: "subject-value",
					Body:    "body-value",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkForEmpties(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("checkForEmpties() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isValid(t *testing.T) {
	type args struct {
		p        *PostRequestBody
		protocol string
		network  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"err-nil",
			args{},
			true,
		},
		{
			"err-empties",
			args{
				&PostRequestBody{},
				"ethereum",
				"mainnet",
			},
			true,
		},
		{
			"err-parse-to",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers:   &PostHeaders{},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "public-key-value",
					},
				},
				"ethereum",
				"mainnet",
			},
			true,
		},
		{
			"err-parse-from",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To: hex.EncodeToString(addresstest.EthereumCharlotte),
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "public-key-value",
					},
				},
				"ethereum",
				"mainnet",
			},
			true,
		},
		{
			"err-reply-to",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To:      hex.EncodeToString(addresstest.EthereumCharlotte),
							From:    hex.EncodeToString(addresstest.EthereumSofia),
							ReplyTo: "<invalid",
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "public-key-value",
					},
				},
				"ethereum",
				"mainnet",
			},
			true,
		},
		{
			"err-public-key",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To:   hex.EncodeToString(addresstest.EthereumCharlotte),
							From: hex.EncodeToString(addresstest.EthereumSofia),
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "public-key-value",
					},
				},
				"ethereum",
				"mainnet",
			},
			true,
		},
		{
			"err-encryption-method-name",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To:   hex.EncodeToString(addresstest.EthereumCharlotte),
							From: hex.EncodeToString(addresstest.EthereumSofia),
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: hex.EncodeToString(addresstest.EthereumSofia),
					},
				},
				"ethereum",
				"mainnet",
			},
			true,
		},
		// these tests should be brought back in some form
		// {
		// 	"err-decode-to",
		// 	args{
		// 		&PostRequestBody{
		// 			Message: PostMessage{
		// 				Headers: &PostHeaders{
		// 					To:   hex.EncodeToString(addresstest.EthereumCharlotte),
		// 					From: hex.EncodeToString(addresstest.EthereumSofia),
		// 				},
		// 				Subject:   "subject-value",
		// 				Body:      "body-value",
		// 				PublicKey: "0x" + hex.EncodeToString(testutil.CharlottePublicKey.Bytes()),
		// 			},
		// 		},
		// 		"ethereum",
		// 		"mainnet",
		// 	},
		// 	true,
		// },
		// {
		// 	"err-address-from-public-key",
		// 	args{
		// 		&PostRequestBody{
		// 			Message: PostMessage{
		// 				Headers: &PostHeaders{
		// 					To:   "0x" + hex.EncodeToString(addresstest.EthereumCharlotte),
		// 					From: hex.EncodeToString(addresstest.EthereumSofia),
		// 				},
		// 				Subject:   "subject-value",
		// 				Body:      "body-value",
		// 				PublicKey: "0x" + hex.EncodeToString(testutil.SofiaPublicKey.Bytes()),
		// 			},
		// 		},
		// 		"ethereum",
		// 		"mainnet",
		// 	},
		// 	true,
		// },
		{
			"success",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To:   "0x" + hex.EncodeToString(addresstest.EthereumCharlotte),
							From: hex.EncodeToString(addresstest.EthereumSofia),
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "0x" + hex.EncodeToString(secp256k1test.CharlottePublicKey.Bytes()),
					},
					EncryptionName: "aes256cbc",
				},
				"ethereum",
				"mainnet",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isValid(tt.args.p, tt.args.protocol, tt.args.network); (err != nil) != tt.wantErr {
				t.Errorf("isValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parsePostRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			args{
				func() *http.Request {
					req := httptest.NewRequest("POST", "/?protocol=ethereum&network=mainnet", strings.NewReader(`
					{
						"message": {
							"body": "test",
							"headers": {
								"from": "0xd5ab4ce3605cd590db609b6b5c8901fdb2ef7fe6",
								"to": "0x92d8f10248c6a3953cc3692a894655ad05d61efb"
							},
							"public-key": "0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d",
							"subject": "test"
						},
						"encryption-method-name": "aes256cbc"
					}
					`))
					return req
				}(),
			},
			false,
			false,
		},
		{
			"err-protocol",
			args{
				func() *http.Request {
					req := httptest.NewRequest("POST", "/?", strings.NewReader(`
					{
						"message": {
							"body": "test",
							"headers": {
								"from": "0xd5ab4ce3605cd590db609b6b5c8901fdb2ef7fe6",
								"to": "0x92d8f10248c6a3953cc3692a894655ad05d61efb"
							},
							"public-key": "0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d",
							"subject": "test"
						}
					}
					`))
					return req
				}(),
			},
			true,
			true,
		},
		{
			"err-network",
			args{
				func() *http.Request {
					req := httptest.NewRequest("POST", "/?protocol=ethereum", strings.NewReader(`
					{
						"message": {
							"body": "test",
							"headers": {
								"from": "0xd5ab4ce3605cd590db609b6b5c8901fdb2ef7fe6",
								"to": "0x92d8f10248c6a3953cc3692a894655ad05d61efb"
							},
							"public-key": "0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d",
							"subject": "test"
						}
					}
					`))
					return req
				}(),
			},
			true,
			true,
		},
		{
			"err-parse-body",
			args{
				func() *http.Request {
					req := httptest.NewRequest("POST", "/?protocol=ethereum&network=mainnet", strings.NewReader(`
					{/}
					`))
					return req
				}(),
			},
			true,
			true,
		},
		{
			"err-invalid-body",
			args{
				func() *http.Request {
					req := httptest.NewRequest("POST", "/?protocol=ethereum&network=mainnet", strings.NewReader(`
					{
						"message": {
							"body": "test",
							"headers": {
								"from": "",
								"to": ""
							},
							"public-key": "0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d",
							"subject": "test"
						}
					}
					`))
					return req
				}(),
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePostRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePostRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("parsePostRequest() go t= %v, wantNil %v", err, tt.wantNil)
				return
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		sent             stores.Sent
		senders          map[string]sender.Message
		ks               keystore.Store
		deriveKeyOptions multi.OptionsBuilders
	}
	tests := []struct {
		name             string
		args             args
		req              *http.Request
		expectedResponse string
		expectedStatus   int
	}{
		// {
		// 	"success",
		// 	args{
		// 		func() stores.Sent {
		// 			m := storestest.NewMockSent(mockCtrl)
		// 			m.EXPECT().PutMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("https://address.com/1620b93b29f056fa752a6a4dd92ff9b3a293580190cf31e24c15278e92fb00f0c28b", "1620b93b29f056fa752a6a4dd92ff9b3a293580190cf31e24c15278e92fb00f0c28b", uint64(1), nil)
		// 			return m
		// 		}(),
		// 		func() map[string]sender.Message {
		// 			return map[string]sender.Message{
		// 				"ethereum/mainnet": func() sender.Message {
		// 					m := sendertest.NewMockMessage(mockCtrl)
		// 					return m
		// 				}(),
		// 			}
		// 		}(),
		// 		func() keystore.Store {
		// 			m := keystoretest.NewMockStore(mockCtrl)
		// 			m.EXPECT().HasAddress([]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x01, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}).Return(true)
		// 			m.EXPECT().GetSigner(
		// 				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x01, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
		// 				"ethereum",
		// 				gomock.Any()).Return(
		// 				func() signer.Signer {
		// 					m := signertest.NewMockSigner(mockCtrl)
		// 					return m
		// 				}(), nil)
		// 			return m
		// 		}(),
		// 		multi.OptionsBuilders{},
		// 	},
		// 	func() *http.Request {
		// 		req := httptest.NewRequest("POST", "/", strings.NewReader(`
		// 		{
		// 			"message": {
		// 				"body": "test",
		// 				"headers": {
		// 					"from": "0xd5ab4ce3605cd590db609b6b5c8901fdb2ef7fe6",
		// 					"to": "0x92d8f10248c6a3953cc3692a894655ad05d61efb"
		// 				},
		// 				"public-key": "0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d",
		// 				"subject": "test"
		// 			}
		// 		}
		// 		`))
		// 		req = mux.SetURLVars(req, map[string]string{
		// 			"network": "mainnet",
		// 		})
		// 		return req
		// 	}(),
		// 	"",
		// 	200,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SendMessage(tt.args.sent, tt.args.senders, tt.args.ks, tt.args.deriveKeyOptions))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.expectedStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}
			if !assert.Equal(tt.expectedResponse, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedResponse)
			}
		})
	}
}
