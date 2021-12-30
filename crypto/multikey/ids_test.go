package multikey

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
)

func TestIDFromPublicKey(t *testing.T) {
	type args struct {
		key crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    byte
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519test.AlicePublicKey,
			},
			0xe2,
			false,
		},
		{
			"secp256k1",
			args{
				secp256k1test.AlicePublicKey,
			},
			0xe1,
			false,
		},
		{
			"sr25519",
			args{
				sr25519test.AlicePublicKey,
			},
			0xe3,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IDFromPublicKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("IDFromPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IDFromPublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIDFromPrivateKey(t *testing.T) {
	type args struct {
		key crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    byte
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519test.AlicePrivateKey,
			},
			0xe2,
			false,
		},
		{
			"secp256k1",
			args{
				secp256k1test.AlicePrivateKey,
			},
			0xe1,
			false,
		},
		{
			"sr25519",
			args{
				sr25519test.AlicePrivateKey,
			},
			0xe3,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IDFromPrivateKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("IDFromPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IDFromPrivateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
