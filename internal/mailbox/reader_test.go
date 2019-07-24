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

package mailbox_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mail/rfc2822"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReadMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	assert := assert.New(t)
	type args struct {
		txData    []byte
		decrypter cipher.Decrypter
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
				testutil.MustHexDecodeString("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					decrypted, _ := ioutil.ReadFile("./testdata/simple.golden.eml")
					m.EXPECT().Decrypt(gomock.Any()).Return(decrypted, nil)
					return m
				}(),
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
				testutil.MustHexDecodeString("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					decrypted, _ := ioutil.ReadFile("./testdata/alternative.golden.eml")
					m.EXPECT().Decrypt(gomock.Any()).Return(decrypted, nil)
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-msg-decrypt",
			args{
				testutil.MustHexDecodeString("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-get-message",
			args{
				testutil.MustHexDecodeString("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("file://TestReadMessage/no_message_at_location-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("file://TestReadMessage/no_message_at_location-2204f3d89e5a"), nil)
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-get-integrity-hash",
			args{
				testutil.MustHexDecodeString("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return([]byte("file://TestReadMessage/no_message_at_location-2204f3d89e5a"), nil)
					m.EXPECT().Decrypt(gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-get-url",
			args{
				testutil.MustHexDecodeString("500801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					m.EXPECT().Decrypt(gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-invalid-envelope",
			args{
				testutil.MustHexDecodeString("000801120f7365637265742d6c6f636174696f6e1a22162032343414ff8b90c8c20d4971b4360d88338bc13e3beb9d4a232adbb5acd67795"),
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					return m
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mailbox.ReadMessage(tt.args.txData, tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("ReadMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
