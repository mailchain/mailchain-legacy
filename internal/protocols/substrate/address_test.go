package substrate

import (
	"testing"

	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_prefixWithNetwork(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		network   string
		publicKey []byte
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
				testutil.MustHexDecodeString("b14d4c84eedf30aabd53ae71286b392f1caaf77597c5a525ea7cc856f91de4a3"),
			},
			[]byte{0x2a, 0xb1, 0x4d, 0x4c, 0x84, 0xee, 0xdf, 0x30, 0xaa, 0xbd, 0x53, 0xae, 0x71, 0x28, 0x6b, 0x39, 0x2f, 0x1c, 0xaa, 0xf7, 0x75, 0x97, 0xc5, 0xa5, 0x25, 0xea, 0x7c, 0xc8, 0x56, 0xf9, 0x1d, 0xe4, 0xa3},
			false,
		},
		{
			"invalid",
			args{
				"invalid",
				testutil.MustHexDecodeString("b14d4c84eedf30aabd53ae71286b392f1caaf77597c5a525ea7cc856f91de4a3"),
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
		publicKey []byte
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
				testutil.MustHexDecodeString("b14d4c84eedf30aabd53ae71286b392f1caaf77597c5a525ea7cc856f91de4a3"),
			},
			[]byte{0x2a, 0xb1, 0x4d, 0x4c, 0x84, 0xee, 0xdf, 0x30, 0xaa, 0xbd, 0x53, 0xae, 0x71, 0x28, 0x6b, 0x39, 0x2f, 0x1c, 0xaa, 0xf7, 0x75, 0x97, 0xc5, 0xa5, 0x25, 0xea, 0x7c, 0xc8, 0x56, 0xf9, 0x1d, 0xe4, 0xa3, 0x83, 0x20},
			false,
		},
		{
			"err-network",
			args{
				"invalid",
				testutil.MustHexDecodeString("b14d4c84eedf30aabd53ae71286b392f1caaf77597c5a525ea7cc856f91de4a3"),
			},
			nil,
			true,
		},
		{
			"err-key-length",
			args{
				"edgeware-testnet",
				testutil.MustHexDecodeString("b14d"),
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
