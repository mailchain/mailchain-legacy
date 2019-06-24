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

	"github.com/gogo/protobuf/proto"
	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/internal/mail/rfc2822"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/mailbox/signer/signertest"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/sendertest"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/storestest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_defaultEncryptLocation(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		pk       crypto.PublicKey
		location string
	}
	type val struct {
		encryptedLen   int
		pk             crypto.PrivateKey
		wantDecryptErr bool
		location       string
	}
	tests := []struct {
		name    string
		args    args
		val     val
		wantErr bool
	}{
		{
			"testutil-charlotte",
			args{
				testutil.CharlottePublicKey,
				"http://test.com/location",
			},
			val{
				114,
				testutil.CharlottePrivateKey,
				false,
				"http://test.com/location",
			},
			false,
		},
		{
			"testutil-charlotte-incorrect-private-key",
			args{
				testutil.SofiaPublicKey,
				"http://test.com/location",
			},
			val{
				114,
				testutil.CharlottePrivateKey,
				true,
				"",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultEncryptLocation(tt.args.pk, tt.args.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.val.encryptedLen, len(got)) {
				t.Errorf("len(encryptLocation()) = %v, want %v", got, tt.val.encryptedLen)
			}
			if err == nil {
				decrypter := aes256cbc.NewDecrypter(tt.val.pk)
				loc, err := decrypter.Decrypt(got)
				if (err != nil) != tt.val.wantDecryptErr {
					t.Errorf("decrypter.Decrypt() error = %v, wantDecryptErr %v", err, tt.val.wantDecryptErr)
					return
				}
				if !assert.Equal(tt.val.location, string(loc)) {
					t.Errorf("decryptedLocation = %v, want %v", string(loc), tt.val.location)
				}
			}
		})
	}
}

func Test_defaultEncryptMailMessage(t *testing.T) {
	encodedMsg, err := rfc2822.EncodeNewMessage(&mail.Message{
		Headers: &mail.Headers{
			From: mail.Address{DisplayName: "From Display Name", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
			To:   mail.Address{DisplayName: "To Display Name", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
			// Date: time.
			Date:    time.Date(2018, 01, 02, 03, 04, 05, 06, time.UTC),
			Subject: "test",
		},
		ID:   []byte("2c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471002c47eca011e32b52c71005ad8a8f75e1b44c9@mailchain"),
		Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet."),
	})
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)
	type args struct {
		pk         crypto.PublicKey
		encodedMsg []byte
	}
	type val struct {
		encryptedLen   int
		pk             crypto.PrivateKey
		wantDecryptErr bool
		encodedMsg     []byte
	}
	tests := []struct {
		name    string
		args    args
		val     val
		wantErr bool
	}{
		{
			"testutil-charlotte",
			args{
				testutil.CharlottePublicKey,
				encodedMsg,
			},
			val{
				914,
				testutil.CharlottePrivateKey,
				false,
				encodedMsg,
			},
			false,
		},
		{
			"testutil-charlotte-incorrect-private-key",
			args{
				testutil.SofiaPublicKey,
				encodedMsg,
			},
			val{
				914,
				testutil.CharlottePrivateKey,
				true,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultEncryptMailMessage(tt.args.pk, tt.args.encodedMsg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptMailMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.val.encryptedLen, len(got)) {
				t.Errorf("len(encryptMailMessage()) = %v, want %v", got, tt.val.encryptedLen)
			}

			if err == nil {
				decrypter := aes256cbc.NewDecrypter(tt.val.pk)
				m, err := decrypter.Decrypt(got)
				if (err != nil) != tt.val.wantDecryptErr {
					t.Errorf("decrypter.Decrypt() error = %v, wantDecryptErr %v", err, tt.val.wantDecryptErr)
					return
				}
				if !assert.EqualValues(tt.val.encodedMsg, m) {
					t.Errorf("decryptMailMessage = %v, want %v", string(m), tt.val.encodedMsg)
				}
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
	}{
		{
			"success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SendMessage(); !assert.NotNil(got) {
				t.Errorf("Message() = %v", got)
			}
		})
	}
}

func Test_defaultPrefixedBytes(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		data proto.Message
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success",
			args{
				data: &mail.Data{
					EncryptedLocation: []byte("encrypted-location"),
					Hash:              testutil.MustHexDecodeString("1620671f6f840e08b9c6b3e2125e0381dd5da5578a698eb97a357f1015552263aec6"),
				},
			},
			[]byte{0x50, 0x12, 0x12, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x2d, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x22, 0x16, 0x20, 0x67, 0x1f, 0x6f, 0x84, 0xe, 0x8, 0xb9, 0xc6, 0xb3, 0xe2, 0x12, 0x5e, 0x3, 0x81, 0xdd, 0x5d, 0xa5, 0x57, 0x8a, 0x69, 0x8e, 0xb9, 0x7a, 0x35, 0x7f, 0x10, 0x15, 0x55, 0x22, 0x63, 0xae, 0xc6},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultPrefixedBytes(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("prefixedBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("prefixedBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Message(t *testing.T) {
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
		encryptLocation    func(pk crypto.PublicKey, location string) ([]byte, error)
		encryptMailMessage func(pk crypto.PublicKey, encodedMsg []byte) ([]byte, error)
		prefixedBytes      func(data proto.Message) ([]byte, error)
	}
	type funcArgs struct {
		ctx     context.Context
		network string
		msg     *mail.Message
		pubkey  crypto.PublicKey
		sender  sender.Message
		sent    stores.Sent
		signer  signer.Signer
	}
	tests := []struct {
		name     string
		args     args
		funcArgs funcArgs
		wantErr  bool
	}{
		{
			"success",
			args{
				func(pk crypto.PublicKey, location string) ([]byte, error) {
					return []byte("encrypted-location"), nil
				},
				func(pk crypto.PublicKey, encodedMsg []byte) ([]byte, error) {
					return []byte("encrypted-message"), nil
				},
				func(data proto.Message) ([]byte, error) {
					return []byte("prefixed-bytes"), nil
				},
			},
			funcArgs{
				context.Background(),
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					m.EXPECT().Send(gomock.Any(),ethereum.Mainnet, testutil.MustHexDecodeString(msg.Headers.To.ChainAddress), testutil.MustHexDecodeString(msg.Headers.From.ChainAddress), append(encoding.DataPrefix(), []byte("prefixed-bytes")...), signertest.NewMockSigner(mockCtrl), nil).Return(nil)
					return m
				}(),
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					sent.EXPECT().PutMessage(msg.ID, []byte("encrypted-message"), nil).Return("https://location-of-file", nil)
					return sent
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
			},
			false,
		},
		{
			"err-encode-message",
			args{
				func(pk crypto.PublicKey, location string) ([]byte, error) {
					return []byte("encrypted-location"), nil
				},
				func(pk crypto.PublicKey, encodedMsg []byte) ([]byte, error) {
					return []byte("encrypted-message"), nil
				},
				func(data proto.Message) ([]byte, error) {
					return []byte("prefixed-bytes"), nil
				},
			},
			funcArgs{
				context.Background(),
				ethereum.Mainnet,
				nil,
				testutil.CharlottePublicKey,
				func() sender.Message {
					sender := sendertest.NewMockMessage(mockCtrl)
					return sender
				}(),
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					return sent
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
			},
			true,
		},
		{
			"err-encrypt-message",
			args{
				func(pk crypto.PublicKey, location string) ([]byte, error) {
					return []byte("encrypted-location"), nil
				},
				func(pk crypto.PublicKey, encodedMsg []byte) ([]byte, error) {
					return nil, errors.Errorf("encryption failed")
				},
				func(data proto.Message) ([]byte, error) {
					return []byte("prefixed-bytes"), nil
				},
			},
			funcArgs{
				context.Background(),
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() sender.Message {
					sender := sendertest.NewMockMessage(mockCtrl)
					return sender
				}(),
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					return sent
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
			},
			true,
		},
		{
			"err-put-message",
			args{
				func(pk crypto.PublicKey, location string) ([]byte, error) {
					return []byte("encrypted-location"), nil
				},
				func(pk crypto.PublicKey, encodedMsg []byte) ([]byte, error) {
					return []byte("encrypted-message"), nil
				},
				func(data proto.Message) ([]byte, error) {
					return []byte("prefixed-bytes"), nil
				},
			},
			funcArgs{
				context.Background(),
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() sender.Message {
					sender := sendertest.NewMockMessage(mockCtrl)
					return sender
				}(),
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					sent.EXPECT().PutMessage(msg.ID, []byte("encrypted-message"), nil).Return("", errors.Errorf("failed to put message"))
					return sent
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
			},
			true,
		},
		{
			"err-encrypt-location",
			args{
				func(pk crypto.PublicKey, location string) ([]byte, error) {
					return nil, errors.Errorf("failed encrypt location")
				},
				func(pk crypto.PublicKey, encodedMsg []byte) ([]byte, error) {
					return []byte("encrypted-message"), nil
				},
				func(data proto.Message) ([]byte, error) {
					return []byte("prefixed-bytes"), nil
				},
			},
			funcArgs{
				context.Background(),
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() sender.Message {
					sender := sendertest.NewMockMessage(mockCtrl)
					return sender
				}(),
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					sent.EXPECT().PutMessage(msg.ID, []byte("encrypted-message"), nil).Return("https://location-of-file", nil)
					return sent
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
			},
			true,
		},
		{
			"err-prefix",
			args{
				func(pk crypto.PublicKey, location string) ([]byte, error) {
					return []byte("encrypted-location"), nil
				},
				func(pk crypto.PublicKey, encodedMsg []byte) ([]byte, error) {
					return []byte("encrypted-message"), nil
				},
				func(data proto.Message) ([]byte, error) {
					return nil, errors.Errorf("prefix failed")
				},
			},
			funcArgs{
				context.Background(),
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() sender.Message {
					sender := sendertest.NewMockMessage(mockCtrl)
					return sender
				}(),
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					sent.EXPECT().PutMessage(msg.ID, []byte("encrypted-message"), nil).Return("https://location-of-file", nil)
					return sent
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
			},
			true,
		},
		{
			"err-sender",
			args{
				func(pk crypto.PublicKey, location string) ([]byte, error) {
					return []byte("encrypted-location"), nil
				},
				func(pk crypto.PublicKey, encodedMsg []byte) ([]byte, error) {
					return []byte("encrypted-message"), nil
				},
				func(data proto.Message) ([]byte, error) {
					return []byte("prefixed-bytes"), nil
				},
			},
			funcArgs{
				context.Background(),
				ethereum.Mainnet,
				msg,
				testutil.CharlottePublicKey,
				func() sender.Message {
					m := sendertest.NewMockMessage(mockCtrl)
					m.EXPECT().Send(gomock.Any(), ethereum.Mainnet, testutil.MustHexDecodeString(msg.Headers.To.ChainAddress), testutil.MustHexDecodeString(msg.Headers.From.ChainAddress), append(encoding.DataPrefix(), []byte("prefixed-bytes")...), signertest.NewMockSigner(mockCtrl), nil).Return(errors.Errorf("failed sender"))
					return m
				}(),
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					sent.EXPECT().PutMessage(msg.ID, []byte("encrypted-message"), nil).Return("https://location-of-file", nil)
					return sent
				}(),
				func() signer.Signer {
					signer := signertest.NewMockSigner(mockCtrl)
					return signer
				}(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFunc := sendMessage(tt.args.encryptLocation, tt.args.encryptMailMessage, tt.args.prefixedBytes)
			err := gotFunc(tt.funcArgs.ctx, tt.funcArgs.network, tt.funcArgs.msg, tt.funcArgs.pubkey, tt.funcArgs.sender, tt.funcArgs.sent, tt.funcArgs.signer)

			if (err != nil) != tt.wantErr {
				t.Errorf("gotFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
