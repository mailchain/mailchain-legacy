package noop

import (
	"crypto"
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto/cipher"
)

func TestNewEncrypter(t *testing.T) {
	tests := []struct {
		name string
		want Encrypter
	}{
		{
			"success",
			Encrypter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEncrypter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypter_Encrypt(t *testing.T) {
	type args struct {
		recipientPublicKey crypto.PublicKey
		message            cipher.PlainContent
	}
	tests := []struct {
		name    string
		e       Encrypter
		args    args
		want    cipher.EncryptedContent
		wantErr bool
	}{
		{
			"success",
			NewEncrypter(),
			args{
				nil,
				cipher.PlainContent([]byte("test content")),
			},
			bytesEncode(cipher.EncryptedContent([]byte("test content"))),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Encrypter{}
			got, err := e.Encrypt(tt.args.recipientPublicKey, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypter.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encrypter.Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
