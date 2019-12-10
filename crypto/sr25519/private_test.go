package sr25519

import (
	"reflect"
	"testing"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/crypto"
	"github.com/stretchr/testify/assert"
)

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
			"success-sofia-seed",
			args{
				sofiaSeed,
			},
			&sofiaPrivateKey,
			false,
		},
		{
			"success-sofia-bytes",
			args{
				sofiaPrivateKeyBytes,
			},
			&sofiaPrivateKey,
			false,
		},
		{
			"success-charlotte-seed",
			args{
				charlotteSeed,
			},
			&charlottePrivateKey,
			false,
		},
		{
			"success-charlotte-bytes",
			args{
				charlottePrivateKeyBytes,
			},
			&charlottePrivateKey,
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

func TestPrivateKey_Kind(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"success",
			"sr25519",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrivateKey{}
			if got := pk.Kind(); got != tt.want {
				t.Errorf("PrivateKey.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_Bytes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		pk   PrivateKey
		want []byte
	}{
		{
			"sucess-sofia",
			sofiaPrivateKey,
			sofiaPrivateKeyBytes,
		},
		{
			"sucess-charllotte",
			charlottePrivateKey,
			charlottePrivateKeyBytes,
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

func TestPrivateKey_PublicKey(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		pk   PrivateKey
		want crypto.PublicKey
	}{
		{
			"sofia",
			sofiaPrivateKey,
			sofiaPublicKey,
		},
		{
			"charlotte",
			charlottePrivateKey,
			charlottePublicKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.PublicKey(); !assert.Equal(tt.want.Bytes(), got.Bytes()) {
				t.Errorf("PrivateKey.PublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_Sign(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		pk      PrivateKey
		msg     []byte
		want    []byte
		wantErr bool
	}{
		{
			"success-charlotte",
			charlottePrivateKey,
			[]byte("message"),
			[]byte{0xc2, 0x9e, 0xeb, 0x12, 0xf1, 0x70, 0x63, 0x55, 0xa2, 0xf, 0x46, 0xcd, 0xe4, 0xf9, 0x21, 0x3d, 0x4a, 0xba, 0x38, 0xb8, 0xb8, 0x41, 0x45, 0x2c, 0x5d, 0x82, 0x6c, 0x48, 0x33, 0xc7, 0x6e, 0x57, 0xad, 0x6c, 0xeb, 0x7b, 0x6a, 0xff, 0x77, 0xdc, 0x48, 0xc4, 0x53, 0x3, 0x77, 0xba, 0xec, 0xc7, 0x71, 0xbd, 0x6d, 0xd3, 0x1e, 0x3e, 0xd5, 0x96, 0x49, 0x58, 0xa3, 0x64, 0x61, 0xc8, 0x6d, 0x89},
			false,
		},
		{
			"success-sofia",
			sofiaPrivateKey,
			[]byte("egassem"),
			[]byte{0x30, 0x5e, 0xaa, 0x1c, 0x13, 0x1c, 0x89, 0xb4, 0x85, 0xbd, 0x6d, 0x8b, 0x48, 0xe0, 0xf6, 0x8, 0xb, 0x29, 0xf3, 0x5b, 0x5f, 0xf, 0xca, 0x8c, 0x36, 0x47, 0x7a, 0xf2, 0x9, 0x83, 0x1d, 0x79, 0x33, 0x5, 0x8b, 0x95, 0xcd, 0x62, 0xe2, 0x5b, 0x9f, 0x91, 0xcb, 0xc2, 0x60, 0x75, 0xe8, 0x7d, 0xe9, 0xe6, 0xbc, 0x4, 0xbf, 0x16, 0x5a, 0x39, 0x35, 0x9a, 0xb8, 0x53, 0xde, 0x41, 0xbc, 0x8a},
			false,
		},
		{
			"err-len",
			PrivateKey{key: func() *schnorrkel.SecretKey {
				priv, err := keyFromSeed([]byte{0xd, 0x9b, 0x4a, 0x3c, 0x10, 0x72, 0x19, 0x91, 0xc6, 0xb8, 0x6, 0xf0, 0xf3, 0x43, 0x53, 0x5d, 0xc2, 0xb4, 0x6c, 0x74})
				if err != nil {
					return nil
				}
				return priv
			}()},
			[]byte("message"),
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pk.Sign(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(len(tt.want), len(got)) {
				t.Errorf("Sign() = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_keyFromBytes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		input   []byte
		want    *PrivateKey
		wantErr bool
	}{
		{
			"error-32bytes-key",
			sofiaPublicKeyBytes,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keyFromBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("keyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("keyFromBytes() = %v,\n want %v", got, tt.want)
			}
		})
	}
}
