package ed25519

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKey_Bytes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		pk   PrivateKey
		want []byte
	}{
		{
			"sofia",
			func() PrivateKey {
				v := sofiaPrivateKey()
				return *v
			}(),
			sofiaPrivateKeyBytes(),
		},
		{
			"charlotte",
			func() PrivateKey {
				v := charlottePrivateKey()
				return *v
			}(),
			charlottePrivateKeyBytes(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.Bytes(); !assert.Equal(tt.want, got) {
				t.Errorf("PrivateKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKeyFromBytes(t *testing.T) {
	type args struct {
		pk []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *PrivateKey
		wantErr bool
	}{
		{
			"success-sofia",
			args{
				sofiaSeed(),
			},
			sofiaPrivateKey(),
			false,
		},
		{
			"success-charlotte",
			args{
				charlotteSeed(),
			},
			charlottePrivateKey(),
			false,
		},
		{
			"err-len",
			args{
				testutil.MustHexDecodeString("39d4c9"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyFromBytes(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_PublicKey(t *testing.T) {
	tests := []struct {
		name string
		pk   PrivateKey
		want crypto.PublicKey
	}{
		{
			"sofia",
			func() PrivateKey {
				v := sofiaPrivateKey()
				return *v
			}(),
			sofiaPublicKey(),
		},
		{
			"charlotte",
			func() PrivateKey {
				v := charlottePrivateKey()
				return *v
			}(),
			charlottePublicKey(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.PublicKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKey.PublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
