package nacl

import (
	"crypto/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ecdh"
	"github.com/mailchain/mailchain/crypto/cryptotest"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/stretchr/testify/assert"
)

func TestNewDecrypter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		privateKey crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    *Decrypter
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519test.BobPrivateKey,
			},
			&Decrypter{
				privateKey: ed25519test.BobPrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewED25519(rand.Reader)
					return k
				}(),
			},
			false,
		},
		{
			"sr25519",
			args{
				sr25519test.BobPrivateKey,
			},
			&Decrypter{
				privateKey: sr25519test.BobPrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSR25519(rand.Reader)
					return k
				}(),
			},
			false,
		},
		{
			"secp256k1",
			args{
				secp256k1test.BobPrivateKey,
			},
			&Decrypter{
				privateKey: secp256k1test.BobPrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSECP256K1(rand.Reader)
					return k
				}(),
			},
			false,
		}, {
			"err-invalid-key",
			args{
				cryptotest.NewMockPrivateKey(mockCtrl),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDecrypter(tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDecrypter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("NewDecrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypter_Decrypt(t *testing.T) {
	type fields struct {
		privateKey  crypto.PrivateKey
		keyExchange cipher.KeyExchange
	}
	type args struct {
		data cipher.EncryptedContent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    cipher.PlainContent
		wantErr bool
	}{
		{
			"ed25519-bob",
			fields{
				privateKey: ed25519test.BobPrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewED25519(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe2, 0x80, 0x56, 0xec, 0xbf, 0x3c, 0xc5, 0xac, 0xd1, 0x60, 0xdd, 0xf0, 0x22, 0x97, 0xbb, 0xba, 0xa1, 0x55, 0x5b, 0xde, 0xa0, 0x4, 0xc2, 0x9b, 0xa9, 0x96, 0x48, 0x89, 0xe1, 0xdc, 0xcd, 0x1, 0x5b, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0xc2, 0x45, 0x87, 0x8, 0x51, 0xe3, 0x9b, 0xff, 0x31, 0x8, 0x9f, 0x40, 0xb6, 0x99, 0x57, 0x99, 0xef, 0x20, 0xba, 0x8e, 0x3b, 0xfd, 0xd1},
			},
			cipher.PlainContent{0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			false,
		},
		{
			"ed25519-alice",
			fields{
				privateKey: ed25519test.AlicePrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewED25519(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe2, 0x80, 0x56, 0xec, 0xbf, 0x3c, 0xc5, 0xac, 0xd1, 0x60, 0xdd, 0xf0, 0x22, 0x97, 0xbb, 0xba, 0xa1, 0x55, 0x5b, 0xde, 0xa0, 0x4, 0xc2, 0x9b, 0xa9, 0x96, 0x48, 0x89, 0xe1, 0xdc, 0xcd, 0x1, 0x5b, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x77, 0x66, 0x5d, 0x95, 0x6d, 0x2f, 0x8e, 0x7, 0x7e, 0x90, 0x7, 0xa4, 0xa1, 0xff, 0x59, 0x9c, 0xbf, 0xf9, 0x38, 0x16, 0xd6, 0x8e, 0xed},
			},
			cipher.PlainContent{0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			false,
		},
		{
			"sr25519-bob",
			fields{
				privateKey: sr25519test.BobPrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSR25519(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe3, 0x82, 0xde, 0x3d, 0x1f, 0x8a, 0x6, 0x59, 0xd2, 0xc6, 0x39, 0xa1, 0x8e, 0x6d, 0x59, 0x3, 0x18, 0x8b, 0x5d, 0xf2, 0x68, 0xc, 0x52, 0x27, 0x61, 0x36, 0x6f, 0xa6, 0xfb, 0x92, 0xde, 0x8a, 0xc, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0xfb, 0x91, 0x80, 0x5f, 0xe7, 0xc7, 0x35, 0xa4, 0x22, 0x50, 0x80, 0x25, 0x53, 0x15, 0xdd, 0x9c, 0x18, 0x93, 0xdc, 0xa3, 0x20, 0xa2, 0x51},
			},
			cipher.PlainContent{0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			false,
		},
		{
			"sr25519-alice",
			fields{
				privateKey: sr25519test.AlicePrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSR25519(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe3, 0x82, 0xde, 0x3d, 0x1f, 0x8a, 0x6, 0x59, 0xd2, 0xc6, 0x39, 0xa1, 0x8e, 0x6d, 0x59, 0x3, 0x18, 0x8b, 0x5d, 0xf2, 0x68, 0xc, 0x52, 0x27, 0x61, 0x36, 0x6f, 0xa6, 0xfb, 0x92, 0xde, 0x8a, 0xc, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x2f, 0xe4, 0xca, 0xb, 0x84, 0x4e, 0x10, 0x29, 0x70, 0x2a, 0xe1, 0xb0, 0xeb, 0x4f, 0x96, 0x2c, 0xfe, 0xaa, 0x73, 0x7f, 0xc9, 0x69, 0x1b},
			},
			cipher.PlainContent{0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			false,
		},
		{
			"secp256k1-bob",
			fields{
				privateKey: secp256k1test.BobPrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSECP256K1(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe1, 0x2, 0xa7, 0xc3, 0xc4, 0xf5, 0x83, 0x73, 0xc0, 0xf6, 0x30, 0xc8, 0x62, 0x63, 0xf, 0x6d, 0x8a, 0xbd, 0xe1, 0x39, 0x48, 0x30, 0xb9, 0xa4, 0x98, 0x8a, 0x3d, 0x6e, 0xe8, 0x86, 0x8b, 0x7a, 0x45, 0xf7, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0xdc, 0xcc, 0x41, 0xa1, 0x3a, 0xa5, 0x54, 0xb7, 0x6b, 0xa3, 0x76, 0x72, 0x21, 0x7d, 0xca, 0xe4, 0xea, 0x43, 0xea, 0xc2, 0x57, 0x83, 0xa9},
			},
			cipher.PlainContent{0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			false,
		},
		{
			"secp256k1-alice",
			fields{
				privateKey: secp256k1test.AlicePrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSECP256K1(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe1, 0x2, 0xa7, 0xc3, 0xc4, 0xf5, 0x83, 0x73, 0xc0, 0xf6, 0x30, 0xc8, 0x62, 0x63, 0xf, 0x6d, 0x8a, 0xbd, 0xe1, 0x39, 0x48, 0x30, 0xb9, 0xa4, 0x98, 0x8a, 0x3d, 0x6e, 0xe8, 0x86, 0x8b, 0x7a, 0x45, 0xf7, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x4c, 0xc6, 0xb4, 0xfe, 0x90, 0xda, 0x5e, 0xc7, 0x42, 0x23, 0x12, 0x5f, 0x6c, 0xb3, 0xf5, 0x15, 0x70, 0xed, 0xa7, 0x78, 0x8b, 0xf, 0x3d},
			},
			cipher.PlainContent{0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			false,
		},
		{
			"err-decode",
			fields{
				privateKey: secp256k1test.AlicePrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSECP256K1(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x00, 0xe1, 0x2, 0xa7, 0xc3, 0xc4, 0xf5, 0x83, 0x73, 0xc0, 0xf6, 0x30, 0xc8, 0x62, 0x63, 0xf, 0x6d, 0x8a, 0xbd, 0xe1, 0x39, 0x48, 0x30, 0xb9, 0xa4, 0x98, 0x8a, 0x3d, 0x6e, 0xe8, 0x86, 0x8b, 0x7a, 0x45, 0xf7, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x4c, 0xc6, 0xb4, 0xfe, 0x90, 0xda, 0x5e, 0xc7, 0x42, 0x23, 0x12, 0x5f, 0x6c, 0xb3, 0xf5, 0x15, 0x70, 0xed, 0xa7, 0x78, 0x8b, 0xf, 0x3d},
			},
			nil,
			true,
		},
		{
			"err-secp256k1-wrong-key",
			fields{
				privateKey: secp256k1test.BobPrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSECP256K1(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe1, 0x2, 0xa7, 0xc3, 0xc4, 0xf5, 0x83, 0x73, 0xc0, 0xf6, 0x30, 0xc8, 0x62, 0x63, 0xf, 0x6d, 0x8a, 0xbd, 0xe1, 0x39, 0x48, 0x30, 0xb9, 0xa4, 0x98, 0x8a, 0x3d, 0x6e, 0xe8, 0x86, 0x8b, 0x7a, 0x45, 0xf7, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x4c, 0xc6, 0xb4, 0xfe, 0x90, 0xda, 0x5e, 0xc7, 0x42, 0x23, 0x12, 0x5f, 0x6c, 0xb3, 0xf5, 0x15, 0x70, 0xed, 0xa7, 0x78, 0x8b, 0xf, 0x3d},
			},
			nil,
			true,
		},
		{
			"err-ed25519-wrong-key",
			fields{
				privateKey: ed25519test.BobPrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewED25519(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe2, 0x80, 0x56, 0xec, 0xbf, 0x3c, 0xc5, 0xac, 0xd1, 0x60, 0xdd, 0xf0, 0x22, 0x97, 0xbb, 0xba, 0xa1, 0x55, 0x5b, 0xde, 0xa0, 0x4, 0xc2, 0x9b, 0xa9, 0x96, 0x48, 0x89, 0xe1, 0xdc, 0xcd, 0x1, 0x5b, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x77, 0x66, 0x5d, 0x95, 0x6d, 0x2f, 0x8e, 0x7, 0x7e, 0x90, 0x7, 0xa4, 0xa1, 0xff, 0x59, 0x9c, 0xbf, 0xf9, 0x38, 0x16, 0xd6, 0x8e, 0xed},
			},
			nil,
			true,
		},
		{
			"err-ed25519-wrong-key",
			fields{
				privateKey: sr25519test.AlicePrivateKey,
				keyExchange: func() cipher.KeyExchange {
					k, _ := ecdh.NewSR25519(rand.Reader)
					return k
				}(),
			},
			args{
				cipher.EncryptedContent{0x2a, 0xe3, 0x82, 0xde, 0x3d, 0x1f, 0x8a, 0x6, 0x59, 0xd2, 0xc6, 0x39, 0xa1, 0x8e, 0x6d, 0x59, 0x3, 0x18, 0x8b, 0x5d, 0xf2, 0x68, 0xc, 0x52, 0x27, 0x61, 0x36, 0x6f, 0xa6, 0xfb, 0x92, 0xde, 0x8a, 0xc, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0xfb, 0x91, 0x80, 0x5f, 0xe7, 0xc7, 0x35, 0xa4, 0x22, 0x50, 0x80, 0x25, 0x53, 0x15, 0xdd, 0x9c, 0x18, 0x93, 0xdc, 0xa3, 0x20, 0xa2, 0x51},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Decrypter{
				privateKey:  tt.fields.privateKey,
				keyExchange: tt.fields.keyExchange,
			}
			got, err := d.Decrypt(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypter.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Decrypter.Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
