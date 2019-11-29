package sr25519

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestPublicKey_Bytes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		pk   PublicKey
		want []byte
	}{
		{
			"sofia",
			sofiaPublicKey,
			sofiaPublicKeyBytes,
		},
		{
			"charlotte",
			charlottePublicKey,
			charlottePublicKeyBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.Bytes(); !assert.Equal(tt.want, got) {
				t.Errorf("PublicKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKey_Kind(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		pk   PublicKey
		want string
	}{
		{
			"charlotte",
			charlottePublicKey,
			crypto.SR25519,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.Kind(); !assert.Equal(tt.want, got) {
				t.Errorf("PublicKey.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKeyFromBytes(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		keyBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *PublicKey
		wantErr bool
	}{
		{
			"success-sofia-bytes",
			args{
				sofiaPublicKeyBytes,
			},
			&sofiaPublicKey,
			false,
		},
		{
			"err-too-short",
			args{
				[]byte{0x72, 0x3c, 0xaa, 0x23},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromBytes(tt.args.keyBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("PublicKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
