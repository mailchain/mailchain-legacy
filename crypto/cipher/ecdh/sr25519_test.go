package ecdh

import (
	"bytes"
	"crypto/rand"
	"io"
	"reflect"
	"testing"
	"testing/iotest"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/stretchr/testify/assert"
)

func TestNewSR25519(t *testing.T) {
	type args struct {
		rand io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *SR25519
		wantErr bool
	}{
		{
			"success",
			args{
				rand.Reader,
			},
			&SR25519{
				rand.Reader,
			},
			false,
		},
		{
			"err-nil-rand",
			args{
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSR25519(tt.args.rand)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSR25519() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSR25519() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSR25519_EphemeralKey(t *testing.T) {
	type fields struct {
		rand io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"success",
			fields{
				rand.Reader,
			},
			false,
		},
		{
			"err-rand",
			fields{
				iotest.DataErrReader(bytes.NewReader(nil)),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kx := SR25519{
				rand: tt.fields.rand,
			}
			_, err := kx.EphemeralKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("SR25519.EphemeralKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSR25519_publicKey(t *testing.T) {
	type fields struct {
		rand io.Reader
	}
	type args struct {
		pubKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantKey crypto.PublicKey
		wantErr bool
	}{
		{
			"success-sr25519-sofia",
			fields{
				nil,
			},
			args{
				sr25519test.SofiaPublicKey,
			},
			sr25519test.SofiaPublicKey,
			false,
		},
		{
			"success-sr25519-charlotte",
			fields{
				nil,
			},
			args{
				sr25519test.CharlottePublicKey,
			},
			sr25519test.CharlottePublicKey,
			false,
		},
		{
			"err-secp256k1-sofia",
			fields{
				nil,
			},
			args{
				secp256k1test.SofiaPublicKey,
			},
			(*sr25519.PublicKey)(nil),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			kx := SR25519{
				rand: tt.fields.rand,
			}
			gotKey, err := kx.publicKey(tt.args.pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SR25519.publicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantKey, gotKey) {
				t.Errorf("SR25519.publicKey() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestSR25519_privateKey(t *testing.T) {
	type fields struct {
		rand io.Reader
	}
	type args struct {
		privKey crypto.PrivateKey
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantKey crypto.PrivateKey
		wantErr bool
	}{
		{
			"success-sr25519-sofia",
			fields{
				nil,
			},
			args{
				sr25519test.SofiaPrivateKey,
			},
			sr25519test.SofiaPrivateKey,
			false,
		},
		{
			"success-ed25519-charlotte",
			fields{
				nil,
			},
			args{
				sr25519test.CharlottePrivateKey,
			},
			sr25519test.CharlottePrivateKey,
			false,
		},
		{
			"err-secp256k1-sofia",
			fields{
				nil,
			},
			args{
				secp256k1test.SofiaPrivateKey,
			},
			(*sr25519.PrivateKey)(nil),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			kx := SR25519{
				rand: tt.fields.rand,
			}
			gotKey, err := kx.privateKey(tt.args.privKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SR25519.privateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantKey, gotKey) {
				t.Errorf("SR25519.privateKey() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestSR25519_SharedSecret(t *testing.T) {
	type fields struct {
		rand io.Reader
	}
	type args struct {
		ephemeralKey crypto.PrivateKey
		recipientKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-charlotte-sofia",
			fields{
				nil,
			},
			args{
				sr25519test.CharlottePrivateKey,
				sr25519test.SofiaPublicKey,
			},
			[]byte{0x59, 0x77, 0xa4, 0xa1, 0xc9, 0x50, 0xe9, 0xa7, 0xfb, 0x91, 0x4b, 0xba, 0x27, 0x2f, 0x6d, 0x81, 0xbb, 0xa7, 0xd2, 0xf6, 0x7b, 0xff, 0x42, 0x55, 0xea, 0xdf, 0xeb, 0x83, 0x8e, 0x48, 0xc1, 0x7a, 0xe0, 0xed, 0xf6, 0xab, 0x45, 0x8d, 0x70, 0x5, 0xb9, 0x63, 0x95, 0x2f, 0xd7, 0xa0, 0xe, 0xb7, 0x7, 0x5e, 0x72, 0x48, 0x1c, 0xaa, 0xdd, 0x95, 0x39, 0x5e, 0x5c, 0x3c, 0x9b, 0x1c, 0xbd, 0x25},
			false,
		},
		{
			"success-sofia-charlotte",
			fields{
				nil,
			},
			args{
				sr25519test.SofiaPrivateKey,
				sr25519test.CharlottePublicKey,
			},
			[]byte{0x59, 0x77, 0xa4, 0xa1, 0xc9, 0x50, 0xe9, 0xa7, 0xfb, 0x91, 0x4b, 0xba, 0x27, 0x2f, 0x6d, 0x81, 0xbb, 0xa7, 0xd2, 0xf6, 0x7b, 0xff, 0x42, 0x55, 0xea, 0xdf, 0xeb, 0x83, 0x8e, 0x48, 0xc1, 0x7a, 0xe0, 0xed, 0xf6, 0xab, 0x45, 0x8d, 0x70, 0x5, 0xb9, 0x63, 0x95, 0x2f, 0xd7, 0xa0, 0xe, 0xb7, 0x7, 0x5e, 0x72, 0x48, 0x1c, 0xaa, 0xdd, 0x95, 0x39, 0x5e, 0x5c, 0x3c, 0x9b, 0x1c, 0xbd, 0x25},
			false,
		},
		{
			"err-sofia-sofia",
			fields{
				nil,
			},
			args{
				sr25519test.SofiaPrivateKey,
				sr25519test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"err-charlotte-charlotte",
			fields{
				nil,
			},
			args{
				sr25519test.CharlottePrivateKey,
				sr25519test.CharlottePublicKey,
			},
			nil,
			true,
		},
		{
			"err-private-key",
			fields{
				nil,
			},
			args{
				nil,
				sr25519test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"err-public-key",
			fields{
				nil,
			},
			args{
				sr25519test.CharlottePrivateKey,
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			kx := SR25519{
				rand: tt.fields.rand,
			}
			got, err := kx.SharedSecret(tt.args.ephemeralKey, tt.args.recipientKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SR25519.SharedSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("SR25519.SharedSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
