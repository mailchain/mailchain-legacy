package ecdh

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"reflect"
	"testing"
	"testing/iotest"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/stretchr/testify/assert"
)

func TestNewSECP256K1(t *testing.T) {
	type args struct {
		rand io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *SECP256K1
		wantErr bool
	}{
		{
			"success",
			args{
				rand.Reader,
			},
			&SECP256K1{
				rand.Reader,
				ethcrypto.S256(),
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
			got, err := NewSECP256K1(tt.args.rand)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSECP256K1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSECP256K1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSECP256K1_EphemeralKey(t *testing.T) {
	type fields struct {
		rand  io.Reader
		curve elliptic.Curve
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
				ethcrypto.S256(),
			},
			false,
		},
		{
			"err-rand",
			fields{
				iotest.DataErrReader(bytes.NewReader(nil)),
				ethcrypto.S256(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kx := SECP256K1{
				rand:  tt.fields.rand,
				curve: tt.fields.curve,
			}
			_, err := kx.EphemeralKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("SECP256K1.EphemeralKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSECP256K1_publicKey(t *testing.T) {
	type fields struct {
		rand  io.Reader
		curve elliptic.Curve
	}
	type args struct {
		pubKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"success-secp256k1-sofia",
			fields{
				nil,
				nil,
			},
			args{
				secp256k1test.SofiaPublicKey,
			},
			false,
		},
		{
			"success-secp256k1-charlotte",
			fields{
				nil,
				nil,
			},
			args{
				secp256k1test.CharlottePublicKey,
			},
			false,
		},
		{
			"err-ed25519-sofia",
			fields{
				nil,
				nil,
			},
			args{
				ed25519test.SofiaPublicKey},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kx := SECP256K1{
				rand:  tt.fields.rand,
				curve: tt.fields.curve,
			}
			_, err := kx.publicKey(tt.args.pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SECP256K1.publicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSECP256K1_privateKey(t *testing.T) {
	type fields struct {
		rand  io.Reader
		curve elliptic.Curve
	}
	type args struct {
		privKey crypto.PrivateKey
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"success-secp256k1-sofia",
			fields{
				nil,
				nil,
			},
			args{
				secp256k1test.SofiaPrivateKey,
			},
			false,
		},
		{
			"success-secp256k1-charlotte",
			fields{
				nil,
				nil,
			},
			args{
				secp256k1test.CharlottePrivateKey,
			},
			false,
		},
		{
			"err-ed25519-sofia",
			fields{
				nil,
				nil,
			},
			args{
				ed25519test.SofiaPrivateKey,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kx := SECP256K1{
				rand:  tt.fields.rand,
				curve: tt.fields.curve,
			}
			_, err := kx.privateKey(tt.args.privKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SECP256K1.privateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSECP256K1_SharedSecret(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		rand  io.Reader
		curve elliptic.Curve
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
				ethcrypto.S256(),
			},
			args{
				secp256k1test.CharlottePrivateKey,
				secp256k1test.SofiaPublicKey,
			},
			[]byte{0xb6, 0xbd, 0xfa, 0xde, 0x23, 0x17, 0x82, 0x72, 0x42, 0x5d, 0x25, 0x77, 0x4a, 0x7d, 0xd, 0x38, 0x8f, 0xbe, 0xf9, 0x48, 0x8, 0x93, 0xfc, 0xc3, 0x64, 0x6a, 0xcc, 0xc1, 0x23, 0xea, 0xcc, 0x47},
			false,
		},
		{
			"success-sofia-charlotte",
			fields{
				nil,
				ethcrypto.S256(),
			},
			args{
				secp256k1test.SofiaPrivateKey,
				secp256k1test.CharlottePublicKey,
			},
			[]byte{0xb6, 0xbd, 0xfa, 0xde, 0x23, 0x17, 0x82, 0x72, 0x42, 0x5d, 0x25, 0x77, 0x4a, 0x7d, 0xd, 0x38, 0x8f, 0xbe, 0xf9, 0x48, 0x8, 0x93, 0xfc, 0xc3, 0x64, 0x6a, 0xcc, 0xc1, 0x23, 0xea, 0xcc, 0x47},
			false,
		},
		{
			"err-sofia-sofia",
			fields{
				nil,
				ethcrypto.S256(),
			},
			args{
				secp256k1test.SofiaPrivateKey,
				secp256k1test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"err-charlotte-charlotte",
			fields{
				nil,
				ethcrypto.S256(),
			},
			args{
				secp256k1test.CharlottePrivateKey,
				secp256k1test.CharlottePublicKey,
			},
			nil,
			true,
		},
		{
			"err-private-key",
			fields{
				nil,
				ethcrypto.S256(),
			},
			args{
				nil,
				secp256k1test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"err-public-key",
			fields{
				nil,
				ethcrypto.S256(),
			},
			args{
				secp256k1test.CharlottePrivateKey,
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kx := SECP256K1{
				rand:  tt.fields.rand,
				curve: tt.fields.curve,
			}
			got, err := kx.SharedSecret(tt.args.ephemeralKey, tt.args.recipientKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SECP256K1.SharedSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("SECP256K1.SharedSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
