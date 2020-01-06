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
	assert := assert.New(t)
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
		wantKey [32]byte
		wantErr bool
	}{
		{
			"success-SR25519-sofia",
			fields{
				nil,
			},
			args{
				sr25519test.SofiaPublicKey,
			},
			[32]uint8{0x22, 0xc9, 0x82, 0x7, 0x9b, 0x62, 0x1a, 0xc6, 0xfd, 0xba, 0x20, 0x73, 0x71, 0x60, 0xcc, 0x91, 0xb3, 0x8f, 0x75, 0x71, 0x69, 0xfd, 0xfb, 0x97, 0xfe, 0xe7, 0x37, 0xe3, 0x7c, 0x69, 0x19, 0x5b},
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
			[32]uint8{0x9c, 0x19, 0x11, 0x65, 0xc0, 0x42, 0x98, 0x6c, 0x26, 0x5f, 0x3d, 0x62, 0x94, 0x3, 0x2a, 0x7a, 0xe, 0x97, 0x64, 0x7a, 0x1a, 0x1b, 0xde, 0x1d, 0x4d, 0xec, 0x7, 0x9, 0xd6, 0x62, 0x2a, 0x41},
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
			[32]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kx := ED25519{
				rand: tt.fields.rand,
			}
			gotKey, err := kx.publicKey(tt.args.pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ED25519.publicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantKey, gotKey) {
				t.Errorf("ED25519.publicKey() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestSR25519_privateKey(t *testing.T) {
	assert := assert.New(t)
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
		wantKey [32]byte
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
			[32]uint8{0x58, 0x9e, 0xe, 0x2a, 0x34, 0x4, 0x8f, 0xb7, 0xa2, 0x3a, 0xe1, 0xa, 0xcb, 0xe0, 0xd3, 0x2b, 0x9b, 0x7f, 0xf7, 0x44, 0x25, 0xc4, 0x80, 0xa8, 0xf7, 0xc2, 0xeb, 0xea, 0xf0, 0xff, 0x77, 0x76},
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
			[32]uint8{0xd8, 0x9, 0x35, 0xbd, 0xce, 0x18, 0xc1, 0x87, 0x54, 0xbe, 0x74, 0x84, 0xf5, 0xbf, 0xa6, 0x1d, 0x87, 0x60, 0xfd, 0xb4, 0x3a, 0x9d, 0x98, 0x86, 0x50, 0x28, 0x22, 0x21, 0x8a, 0xe, 0xc6, 0x6b},
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
			[32]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kx := ED25519{
				rand: tt.fields.rand,
			}
			gotKey, err := kx.privateKey(tt.args.privKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ED25519.privateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantKey, gotKey) {
				t.Errorf("ED25519.privateKey() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestSR25519_SharedSecret(t *testing.T) {
	assert := assert.New(t)
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
			[]byte{0xf1, 0x48, 0xbc, 0xc6, 0xb7, 0x3d, 0x8c, 0xb3, 0xdd, 0x85, 0x8e, 0x26, 0xb1, 0x47, 0x78, 0xfa, 0x9b, 0xfa, 0xc8, 0xc3, 0xdd, 0xd5, 0xdd, 0x9f, 0xe7, 0x1e, 0x26, 0x66, 0xd6, 0x1c, 0xf0, 0x4d},
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
			[]byte{0xf1, 0x48, 0xbc, 0xc6, 0xb7, 0x3d, 0x8c, 0xb3, 0xdd, 0x85, 0x8e, 0x26, 0xb1, 0x47, 0x78, 0xfa, 0x9b, 0xfa, 0xc8, 0xc3, 0xdd, 0xd5, 0xdd, 0x9f, 0xe7, 0x1e, 0x26, 0x66, 0xd6, 0x1c, 0xf0, 0x4d},
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
			if !assert.Equal(tt.want, got) {
				t.Errorf("SR25519.SharedSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
