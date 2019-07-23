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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/mailbox/signer/signertest"
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
			From:    mail.Address{DisplayName: "From Display Name", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
			To:      mail.Address{DisplayName: "To Display Name", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
			Date:    time.Date(2018, 01, 02, 03, 04, 05, 06, time.UTC),
			Subject: "test",
		},
		ID:   []byte("2c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471002c47eca011e32b52c71005ad8a8f75e1b44c9@mailchain"),
		Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet."),
	}

	type args struct {
		ctx          context.Context
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
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(testutil.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)
					m.EXPECT().Encrypt(testutil.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-uint64-bytes"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					m.EXPECT().Send(gomock.Any(), ethereum.Mainnet, testutil.MustHexDecodeString(msg.Headers.To.ChainAddress), testutil.MustHexDecodeString(msg.Headers.From.ChainAddress), gomock.AssignableToTypeOf([]byte{}), signertest.NewMockSigner(mockCtrl), nil).Return(nil)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, []byte{0x16, 0x20, 0xed, 0x2b, 0x9, 0x40, 0x7e, 0xb6, 0x91, 0x67, 0x8a, 0x27, 0x10, 0x42, 0xd2, 0xC3, 0xd2, 0x26, 0xF0, 0xBD, 0x8E, 0x62, 0xF0, 0x7C, 0xF1, 0x61, 0xD2, 0x62, 0x8E, 0x8E, 0x11, 0x7E, 0x35, 0x42}, []byte("encrypted-message"), nil).Return("https://location-of-file", "1620ed2b09407eb691678a271042d2c3d226f0bd8e62f07cf161d2628e8e117e3542", uint64(1), nil)
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
			"err-send",
			args{
				context.Background(),
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(testutil.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)
					m.EXPECT().Encrypt(testutil.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-uint64-bytes"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					m.EXPECT().Send(gomock.Any(), ethereum.Mainnet, testutil.MustHexDecodeString(msg.Headers.To.ChainAddress), testutil.MustHexDecodeString(msg.Headers.From.ChainAddress), gomock.AssignableToTypeOf([]byte{}), signertest.NewMockSigner(mockCtrl), nil).Return(errors.Errorf("failed"))
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, []byte{0x16, 0x20, 0xed, 0x2b, 0x9, 0x40, 0x7e, 0xb6, 0x91, 0x67, 0x8a, 0x27, 0x10, 0x42, 0xd2, 0xC3, 0xd2, 0x26, 0xF0, 0xBD, 0x8E, 0x62, 0xF0, 0x7C, 0xF1, 0x61, 0xD2, 0x62, 0x8E, 0x8E, 0x11, 0x7E, 0x35, 0x42}, []byte("encrypted-message"), nil).Return("https://location-of-file", "1620ed2b09407eb691678a271042d2c3d226f0bd8e62f07cf161d2628e8e117e3542", uint64(1), nil)
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
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(testutil.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, []byte{0x16, 0x20, 0xed, 0x2b, 0x9, 0x40, 0x7e, 0xb6, 0x91, 0x67, 0x8a, 0x27, 0x10, 0x42, 0xd2, 0xC3, 0xd2, 0x26, 0xF0, 0xBD, 0x8E, 0x62, 0xF0, 0x7C, 0xF1, 0x61, 0xD2, 0x62, 0x8E, 0x8E, 0x11, 0x7E, 0x35, 0x42}, []byte("encrypted-message"), nil).Return("https://location-of-file", "1620ed2b09407eb691678a271042d2c3d226f0bd8e62f07cf161d2628e8e117e3542", uint64(1), nil)
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
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(testutil.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, []byte{0x16, 0x20, 0xed, 0x2b, 0x9, 0x40, 0x7e, 0xb6, 0x91, 0x67, 0x8a, 0x27, 0x10, 0x42, 0xd2, 0xC3, 0xd2, 0x26, 0xF0, 0xBD, 0x8E, 0x62, 0xF0, 0x7C, 0xF1, 0x61, 0xD2, 0x62, 0x8E, 0x8E, 0x11, 0x7E, 0x35, 0x42}, []byte("encrypted-message"), nil).Return("https://location-of-file", "1620ed2b09407eb691678a271042d2c3d226f0bd8e62f07cf161d2628e8e117e3542", uint64(255), nil)
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
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(testutil.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return([]byte("encrypted-message"), nil)

					return m
				}(),
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					return m
				}(),
				func() stores.Sent {
					m := storestest.NewMockSent(mockCtrl)
					m.EXPECT().PutMessage(msg.ID, []byte{0x16, 0x20, 0xed, 0x2b, 0x9, 0x40, 0x7e, 0xb6, 0x91, 0x67, 0x8a, 0x27, 0x10, 0x42, 0xd2, 0xC3, 0xd2, 0x26, 0xF0, 0xBD, 0x8E, 0x62, 0xF0, 0x7C, 0xF1, 0x61, 0xD2, 0x62, 0x8E, 0x8E, 0x11, 0x7E, 0x35, 0x42}, []byte("encrypted-message"), nil).Return("", "", uint64(1), errors.Errorf("failed"))
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
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(testutil.CharlottePublicKey, gomock.AssignableToTypeOf([]byte{})).Return(nil, errors.Errorf("failed"))

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
				ethereum.Mainnet,
				nil,
				testutil.CharlottePublicKey,
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
			if err := SendMessage(tt.args.ctx, tt.args.network, tt.args.msg, tt.args.pubkey, tt.args.encrypter, tt.args.msgSender, tt.args.sent, tt.args.msgSigner, tt.args.envelopeKind); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
