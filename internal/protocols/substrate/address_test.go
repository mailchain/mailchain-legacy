package substrate

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_prefixWithNetwork(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		network   string
		publicKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"edgeware-testnet",
			args{
				"edgeware-testnet",
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71},
			false,
		},
		{
			"invalid",
			args{
				"invalid",
				ed25519test.SofiaPublicKey,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prefixWithNetwork(tt.args.network, tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("prefixWithNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("prefixWithNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addSS58Prefix(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		pubKey []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"success",
			args{
				testutil.MustHexDecodeString("b14d"),
			},
			[]byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45, 0xb1, 0x4d},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addSS58Prefix(tt.args.pubKey); !assert.Equal(tt.want, got) {
				t.Errorf("addSS58Prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSS58AddressFormat(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		network   string
		publicKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success",
			args{
				"edgeware-testnet",
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x63, 0x83},
			false,
		},
		{
			"err-network",
			args{
				"invalid",
				ed25519test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"err-key-length",
			args{
				"edgeware-testnet",
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SS58AddressFormat(tt.args.network, tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SS58AddressFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("SS58AddressFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
