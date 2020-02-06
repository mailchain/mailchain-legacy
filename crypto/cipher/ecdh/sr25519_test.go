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
			kx := SR25519{
				rand: tt.fields.rand,
			}
			gotKey, err := kx.publicKey(tt.args.pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SR25519.publicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantKey, gotKey) {
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
			kx := SR25519{
				rand: tt.fields.rand,
			}
			gotKey, err := kx.privateKey(tt.args.privKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SR25519.privateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantKey, gotKey) {
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
			[]byte{0xfe, 0x3d, 0xbc, 0xea, 0xb5, 0x76, 0x69, 0xa8, 0xe6, 0x30, 0x7f, 0xfe, 0x11, 0xf3, 0xac, 0xb9, 0x8f, 0xb9, 0x61, 0xdd, 0x49, 0xef, 0xb1, 0xd9, 0xbc, 0x40, 0x67, 0x7, 0xf7, 0xb0, 0xc8, 0x64},
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
			[]byte{0xfe, 0x3d, 0xbc, 0xea, 0xb5, 0x76, 0x69, 0xa8, 0xe6, 0x30, 0x7f, 0xfe, 0x11, 0xf3, 0xac, 0xb9, 0x8f, 0xb9, 0x61, 0xdd, 0x49, 0xef, 0xb1, 0xd9, 0xbc, 0x40, 0x67, 0x7, 0xf7, 0xb0, 0xc8, 0x64},
			false,
		},
		{
			"success-eve-charlotte",
			fields{
				nil,
			},
			args{
				sr25519test.EvePrivateKey,
				sr25519test.CharlottePublicKey,
			},
			[]byte{0x9f, 0x1f, 0x3a, 0xa8, 0xfc, 0x22, 0xa8, 0x47, 0xed, 0xdd, 0x7a, 0xfc, 0x48, 0x85, 0x80, 0x8f, 0x71, 0x2, 0x12, 0x29, 0xda, 0xf6, 0x9a, 0xb2, 0xba, 0x30, 0x67, 0x76, 0xde, 0x45, 0xff, 0x21},
			false,
		},
		{
			"success-charlotte-eve",
			fields{
				nil,
			},
			args{
				sr25519test.CharlottePrivateKey,
				sr25519test.EvePublicKey,
			},
			[]byte{0x9f, 0x1f, 0x3a, 0xa8, 0xfc, 0x22, 0xa8, 0x47, 0xed, 0xdd, 0x7a, 0xfc, 0x48, 0x85, 0x80, 0x8f, 0x71, 0x2, 0x12, 0x29, 0xda, 0xf6, 0x9a, 0xb2, 0xba, 0x30, 0x67, 0x76, 0xde, 0x45, 0xff, 0x21},
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
			kx := SR25519{
				rand: tt.fields.rand,
			}
			got, err := kx.SharedSecret(tt.args.ephemeralKey, tt.args.recipientKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SR25519.SharedSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("SR25519.SharedSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
