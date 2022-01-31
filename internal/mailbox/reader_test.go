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

package mailbox_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mail/rfc2822"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/storestest"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReadMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		txData    []byte
		decrypter cipher.Decrypter
		cacheFunc func() stores.Cache
	}
	tests := []struct {
		name    string
		args    args
		want    *mail.Message
		wantErr bool
	}{
		{
			"success",
			args{
				encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
				func() cipher.Decrypter {
					decrypted, _ := ioutil.ReadFile("./testdata/simple.golden.eml")
					m := ciphertest.NewMockDecrypter(mockCtrl)

					gomock.InOrder(
						m.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						m.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						m.EXPECT().Decrypt(cipher.EncryptedContent([]byte{0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65})).Return(decrypted, nil),
					)
					return m
				}(),
				func() stores.Cache {
					cache := storestest.NewMockCache(mockCtrl)
					cache.EXPECT().GetMessage(gomock.Any()).Return(nil, errors.New("cache empty"))
					cache.EXPECT().SetMessage(gomock.Any(), gomock.Any())
					return cache
				},
			},
			func() *mail.Message {
				rawMsg, _ := ioutil.ReadFile("./testdata/simple.golden.eml")
				m, _ := rfc2822.DecodeNewMessage(bytes.NewReader(rawMsg))
				return m
			}(),
			false,
		},
		{
			"err-invalid-hash",
			args{
				encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					decrypted, _ := ioutil.ReadFile("./testdata/alternative.golden.eml")
					m.EXPECT().Decrypt(gomock.Any()).Return(decrypted, nil)
					return m
				}(),
				func() stores.Cache {
					cache := storestest.NewMockCache(mockCtrl)
					cache.EXPECT().GetMessage(gomock.Any()).Return(nil, errors.New("cache empty"))
					cache.EXPECT().SetMessage(gomock.Any(), gomock.Any())
					return cache
				},
			},
			nil,
			true,
		},
		{
			"err-msg-decrypt",
			args{
				encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m
				}(),
				func() stores.Cache {
					cache := storestest.NewMockCache(mockCtrl)
					cache.EXPECT().GetMessage(gomock.Any()).Return(nil, errors.New("cache empty"))
					cache.EXPECT().SetMessage(gomock.Any(), gomock.Any())
					return cache
				},
			},
			nil,
			true,
		},
		{
			"err-get-message",
			args{
				encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("file://TestReadMessage/no_message_at_location-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("file://TestReadMessage/no_message_at_location-2204f3d89e5a"), nil)
					return m
				}(),
				func() stores.Cache {
					cache := storestest.NewMockCache(mockCtrl)
					cache.EXPECT().GetMessage(gomock.Any()).Return(nil, errors.New("cache empty"))
					return cache
				},
			},
			nil,
			true,
		},
		{
			"err-get-integrity-hash",
			args{
				encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("file://TestReadMessage/no_message_at_location-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m
				}(),
				func() stores.Cache {
					return storestest.NewMockCache(mockCtrl)
				},
			},
			nil,
			true,
		},
		{
			"err-get-url",
			args{
				encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m
				}(),
				func() stores.Cache {
					return storestest.NewMockCache(mockCtrl)
				},
			},
			nil,
			true,
		},
		{
			"err-invalid-envelope",
			args{
				encodingtest.MustDecodeHex("000801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					return m
				}(),
				func() stores.Cache {
					return storestest.NewMockCache(mockCtrl)
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mailbox.ReadMessage(tt.args.txData, tt.args.decrypter, tt.args.cacheFunc())
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ReadMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
