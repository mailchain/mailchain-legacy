package nacl

import (
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

func Test_getPublicKeyExchange(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		recipientPublicKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    cipher.KeyExchange
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519test.SofiaPublicKey,
			},
			&ecdh.ED25519{},
			false,
		},
		{
			"sr25519",
			args{
				sr25519test.SofiaPublicKey,
			},
			&ecdh.SR25519{},
			false,
		},
		{
			"ed25519",
			args{
				secp256k1test.SofiaPublicKey,
			},
			&ecdh.SECP256K1{},
			false,
		},
		{
			"err-not-supported",
			args{
				cryptotest.NewMockPublicKey(mockCtrl),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPublicKeyExchange(tt.args.recipientPublicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPublicKeyExchange error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(t, tt.want, got) {
				t.Errorf("getPublicKeyExchange = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPrivateKeyExchange(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		recipientPrivateKey crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    cipher.KeyExchange
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519test.SofiaPrivateKey,
			},
			&ecdh.ED25519{},
			false,
		},
		{
			"sr25519",
			args{
				sr25519test.SofiaPrivateKey,
			},
			&ecdh.SR25519{},
			false,
		},
		{
			"ed25519",
			args{
				secp256k1test.SofiaPrivateKey,
			},
			&ecdh.SECP256K1{},
			false,
		},
		{
			"err-not-supported",
			args{
				cryptotest.NewMockPrivateKey(mockCtrl),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPrivateKeyExchange(tt.args.recipientPrivateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPrivateKeyExchange error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(t, tt.want, got) {
				t.Errorf("getPrivateKeyExchange = %v, want %v", got, tt.want)
			}
		})
	}
}
