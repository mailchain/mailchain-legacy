package noop

import (
	"crypto"
	"reflect"
	"testing"

	keys "github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
)

func TestNewEncrypter(t *testing.T) {
	tests := []struct {
		name    string
		want    cipher.Encrypter
		wantErr bool
	}{
		{
			"success",
			&Encrypter{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEncrypter(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypter.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypter_Encrypt(t *testing.T) {
	type fields struct {
		publicKey keys.PublicKey
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
			"success",
			fields{
				nil,
			},
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
			e := Encrypter{publicKey: tt.fields.publicKey}
			got, err := e.Encrypt(tt.args.message)
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
