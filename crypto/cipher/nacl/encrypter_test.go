package nacl

import (
	"bytes"
	"io"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
)

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
		rand io.Reader
	}
	type args struct {
		recipientPublicKey crypto.PublicKey
		message            cipher.PlainContent
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
			},
			args{
				ed25519test.CharlottePublicKey,
				[]byte("message"),
			},
			cipher.EncryptedContent{0x2a, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x5b, 0x19, 0x83, 0xe5, 0x6e, 0x7f, 0xed, 0xfe, 0xbb, 0xd0, 0x70, 0x34, 0xce, 0x25, 0x49, 0x76, 0xa3, 0x50, 0x78, 0x91, 0x18, 0xe6, 0xe3},
			false,
		},
		{
			"err-key-type",
			fields{
				nil,
			},
			args{
				secp256k1.PublicKey{},
				[]byte("message"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Encrypter{
				rand: tt.fields.rand,
			}
			got, err := e.Encrypt(tt.args.recipientPublicKey, tt.args.message)
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
