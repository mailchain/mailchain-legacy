package nacl

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/stretchr/testify/assert"
)

func TestNewEncrypter(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		publicKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    *Encrypter
		wantErr bool
	}{
		{
			"success",
			args{
				ed25519test.CharlottePublicKey,
			},
			&Encrypter{
				rand:      rand.Reader,
				publicKey: ed25519test.CharlottePublicKey,
			},
			false,
		},
		{
			"invalid-key",
			args{
				secp256k1test.CharlottePublicKey,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEncrypter(tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncrypter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(tt.want, got) {
				t.Errorf("NewEncrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validatePublicKeyType(t *testing.T) {
	type args struct {
		recipientPublicKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519.PublicKey{},
			},
			false,
		},
		{
			"sr25519",
			args{
				sr25519.PublicKey{},
			},
			false,
		},
		{
			"not-supported",
			args{
				secp256k1.PublicKey{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePublicKeyType(tt.args.recipientPublicKey); (err != nil) != tt.wantErr {
				t.Errorf("validatePublicKeyType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncrypter_Encrypt(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		rand      io.Reader
		publicKey crypto.PublicKey
	}
	type args struct {
		message cipher.PlainContent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    cipher.EncryptedContent
		wantErr bool
	}{
		{
			"success-charlotte",
			fields{
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
				ed25519test.CharlottePublicKey,
			},
			args{
				[]byte("message"),
			},
			cipher.EncryptedContent{0x2a, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x5b, 0x19, 0x83, 0xe5, 0x6e, 0x7f, 0xed, 0xfe, 0xbb, 0xd0, 0x70, 0x34, 0xce, 0x25, 0x49, 0x76, 0xa3, 0x50, 0x78, 0x91, 0x18, 0xe6, 0xe3},
			false,
		},
		{
			"err-key-type",
			fields{
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
				ed25519test.SofiaPublicKey,
			},
			args{
				[]byte("egassem"),
			},
			cipher.EncryptedContent{0x2a, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0xff, 0xb3, 0x7f, 0x9b, 0x80, 0xe9, 0x85, 0x1f, 0x47, 0xfd, 0xb6, 0xdf, 0x1a, 0x94, 0xc4, 0x7b, 0x92, 0x91, 0x34, 0xf7, 0x76, 0x7e, 0xd4},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Encrypter{
				rand:      tt.fields.rand,
				publicKey: tt.fields.publicKey,
			}
			got, err := e.Encrypt(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypter.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("Encrypter.Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
