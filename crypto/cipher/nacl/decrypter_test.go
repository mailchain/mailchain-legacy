package nacl

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/stretchr/testify/assert"
)

func TestNewDecrypter(t *testing.T) {
	type args struct {
		privateKey crypto.PrivateKey
	}
	tests := []struct {
		name string
		args args
		want Decrypter
	}{
		{
			"success",
			args{
				ed25519test.CharlottePrivateKey,
			},
			Decrypter{
				privateKey: ed25519test.CharlottePrivateKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDecrypter(tt.args.privateKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypter_Decrypt(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		privateKey crypto.PrivateKey
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
			"success-charlotte",
			fields{
				ed25519test.CharlottePrivateKey,
			},
			args{
				cipher.EncryptedContent{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x5b, 0x19, 0x83, 0xe5, 0x6e, 0x7f, 0xed, 0xfe, 0xbb, 0xd0, 0x70, 0x34, 0xce, 0x25, 0x49, 0x76, 0xa3, 0x50, 0x78, 0x91, 0x18, 0xe6, 0xe3},
			},
			cipher.PlainContent{0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Decrypter{
				privateKey: tt.fields.privateKey,
			}
			got, err := d.Decrypt(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypter.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("Decrypter.Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
