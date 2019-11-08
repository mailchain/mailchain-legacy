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

package mailbox

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/mailbox/signer/signertest"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/sendertest"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/storestest"
	"github.com/pkg/errors"
)

func TestSendMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	msg := &mail.Message{
		Headers: &mail.Headers{
			From:        mail.Address{DisplayName: "From Display Name", FullAddress: "0x4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "0x4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
			To:          mail.Address{DisplayName: "To Display Name", FullAddress: "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
			Date:        time.Date(2018, 01, 02, 03, 04, 05, 06, time.UTC),
			Subject:     "test",
			ContentType: "text/plain; charset=\"UTF-8\"",
		},
		ID:   []byte("2c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471002c47eca011e32b52c71005ad8a8f75e1b44c9@mailchain"),
		Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet."),
	}

	type args struct {
		ctx          context.Context
		protocol     string
		network      string
		msg          *mail.Message
		pubkey       crypto.PublicKey
		encrypter    cipher.Encrypter
		msgSender    sender.Message
		sent         stores.Sent
		msgSigner    signer.Signer
		envelopeKind byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				msg,
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-uint64-bytes"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					m.EXPECT().Send(gomock.Any(), ethereum.Mainnet, testutil.MustHexDecodeString(strings.TrimLeft(msg.Headers.To.ChainAddress, "0x")), testutil.MustHexDecodeString(strings.TrimLeft(msg.Headers.From.ChainAddress, "0x")), gomock.AssignableToTypeOf([]byte{}), signertest.NewMockSigner(mockCtrl), nil).Return(nil)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, testutil.MustHexDecodeString("162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d"), []byte("encrypted-message"), nil).Return("https://location-of-file", "162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d", uint64(1), nil)
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				envelope.Kind0x01,
			},
			false,
		},
		{
			"invalid-to",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				&mail.Message{
					Headers: &mail.Headers{
						From:        mail.Address{DisplayName: "From Display Name", FullAddress: "0x4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "0x4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
						To:          mail.Address{DisplayName: "To Display Name", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
						Date:        time.Date(2018, 01, 02, 03, 04, 05, 06, time.UTC),
						Subject:     "test",
						ContentType: "text/plain; charset=\"UTF-8\"",
					},
					ID:   []byte("2c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471002c47eca011e32b52c71005ad8a8f75e1b44c9@mailchain"),
					Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet."),
				},
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-uint64-bytes"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, gomock.AssignableToTypeOf([]byte{}), []byte("encrypted-message"), nil).Return("https://location-of-file", "1620b6a9895ccabf87b802d8507b9fefdd3cf3e4206725f32911bc365caa64cb2248", uint64(1), nil)
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				envelope.Kind0x01,
			},
			true,
		},
		{
			"invalid-from",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				&mail.Message{
					Headers: &mail.Headers{
						From:        mail.Address{DisplayName: "From Display Name", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
						To:          mail.Address{DisplayName: "To Display Name", FullAddress: "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
						Date:        time.Date(2018, 01, 02, 03, 04, 05, 06, time.UTC),
						Subject:     "test",
						ContentType: "text/plain; charset=\"UTF-8\"",
					},
					ID:   []byte("2c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471002c47eca011e32b52c71005ad8a8f75e1b44c9@mailchain"),
					Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet."),
				},
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-uint64-bytes"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, gomock.AssignableToTypeOf([]byte{}), []byte("encrypted-message"), nil).Return("https://location-of-file", "1620dd40e1a3576725761746ea131d08fd8bd49ca94676232295e0f0df00371102f5", uint64(1), nil)
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				envelope.Kind0x01,
			},
			true,
		},
		{
			"err-send",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				msg,
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-uint64-bytes"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					m.EXPECT().Send(gomock.Any(), ethereum.Mainnet, testutil.MustHexDecodeString(strings.TrimLeft(msg.Headers.To.ChainAddress, "0x")), testutil.MustHexDecodeString(strings.TrimLeft(msg.Headers.From.ChainAddress, "0x")), gomock.AssignableToTypeOf([]byte{}), signertest.NewMockSigner(mockCtrl), nil).Return(errors.Errorf("failed"))
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, testutil.MustHexDecodeString("162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d"), []byte("encrypted-message"), nil).Return("https://location-of-file", "162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d", uint64(1), nil)
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				envelope.Kind0x01,
			},
			true,
		},
		{
			"err-new-envelope",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				msg,
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, testutil.MustHexDecodeString("162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d"), []byte("encrypted-message"), nil).Return("https://location-of-file", "162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d", uint64(1), nil)
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				0x00,
			},
			true,
		},
		{
			"err-invalid-loc-code",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				msg,
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, testutil.MustHexDecodeString("162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d"), []byte("encrypted-message"), nil).Return("https://location-of-file", "162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d", uint64(255), nil)
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				0x00,
			},
			true,
		},
		{
			"err-sent-store",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				msg,
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, testutil.MustHexDecodeString("162054f817c0ee9b844de0f294aa23c6cb12cec36a54c1187aaefb06b4a51f39a02d"), []byte("encrypted-message"), nil).Return("", "", uint64(1), errors.Errorf("failed"))
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				0x00,
			},
			true,
		},
		{
			"err-encrypt",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				msg,
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(secp256k1test.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return(nil, errors.Errorf("failed"))

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				0x00,
			},
			true,
		},
		{
			"err-msg-encode",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				nil,
				secp256k1test.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					return m
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
				0x00,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMessage(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.msg, tt.args.pubkey, tt.args.encrypter, tt.args.msgSender, tt.args.sent, tt.args.msgSigner, tt.args.envelopeKind); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
