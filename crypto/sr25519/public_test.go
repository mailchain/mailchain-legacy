package sr25519

import (
	"reflect"
	"testing"
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
