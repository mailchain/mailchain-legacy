package sr25519

import (
	"testing"
	"reflect"

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
	type args struct {
		pk []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *PublicKey
		wantErr bool
	}{
		{
			"success-sofia-seed",
			args{
				sofiaSeed,
			},
			&sofiaPublicKey,
			false,
		},
		{
			"success-sofia-bytes",
			args{
				sofiaPublicKeyBytes,
			},
			&sofiaPublicKey,
			false,
		},
		{
			"success-charlotte-seed",
			args{
				charlotteSeed,
			},
			&charlottePublicKey,
			false,
		},
		{
			"success-charlotte-bytes",
			args{
				charlottePublicKeyBytes,
			},
			&charlottePublicKey,
			false,
		},
		{
			"err-len",
			args{
				[]byte{57, 212, 201},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromBytes(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublicKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
