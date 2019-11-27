package sr25519

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestPublicKey_Bytes(t *testing.T) {
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
			if got := tt.pk.Bytes(); !reflect.DeepEqual(got, tt.want) {
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
