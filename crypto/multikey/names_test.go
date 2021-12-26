package multikey

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
)

func TestKindFromPublicKey(t *testing.T) {
	type args struct {
		key crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519test.AlicePublicKey,
			},
			"ed25519",
			false,
		},
		{
			"secp256k1",
			args{
				secp256k1test.AlicePublicKey,
			},
			"secp256k1",
			false,
		},
		{
			"sr25519",
			args{
				sr25519test.AlicePublicKey,
			},
			"sr25519",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KindFromPublicKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KindFromPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("KindFromPublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKindFromPrivateKey(t *testing.T) {
	type args struct {
		key crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519test.AlicePrivateKey,
			},
			"ed25519",
			false,
		},
		{
			"secp256k1",
			args{
				secp256k1test.AlicePrivateKey,
			},
			"secp256k1",
			false,
		},
		{
			"sr25519",
			args{
				sr25519test.AlicePrivateKey,
			},
			"sr25519",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KindFromPrivateKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KindFromPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("KindFromPrivateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
