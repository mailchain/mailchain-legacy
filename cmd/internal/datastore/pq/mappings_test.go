package pq

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

func TestGetProtocolNetworkUint8(t *testing.T) {
	type args struct {
		protocol string
		network  string
	}
	type result struct {
		protocol uint8
		network  uint8
		wantErr  bool
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{
			"success",
			args{
				protocols.Ethereum,
				ethereum.Mainnet,
			},
			result{
				1,
				1,
				false,
			},
		},
		{
			"err-protocol-unknown",
			args{
				"unknown",
				ethereum.Mainnet,
			},
			result{
				0,
				0,
				true,
			},
		},
		{
			"err-network-unknown",
			args{
				protocols.Ethereum,
				"unknown",
			},
			result{
				0,
				0,
				true,
			},
		},
	}
	for _, tt := range tests {
		uProtocol, uNetwork, err := getProtocolNetworkUint8(tt.args.protocol, tt.args.network)
		if (err != nil) != tt.result.wantErr {
			t.Errorf("getProtocolNetworkUint8() error = %v, wantErr %v", err, tt.result.wantErr)
			return
		}
		assert.Equal(t, tt.result.protocol, uProtocol)
		assert.Equal(t, tt.result.network, uNetwork)
	}
}

func TestGetPublicKeyTypeUint8(t *testing.T) {
	type args struct {
		pub_key_type string
	}
	type result struct {
		pub_key_type uint8
		wantErr      bool
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{
			"success",
			args{
				crypto.KindSECP256K1,
			},
			result{
				crypto.IDSECP256K1,
				false,
			},
		},
		{
			"err-pub_key_type-unknown",
			args{
				"unknown",
			},
			result{
				0,
				true,
			},
		},
	}
	for _, tt := range tests {
		uPubKeyType, err := getPublicKeyTypeUint8(tt.args.pub_key_type)
		if (err != nil) != tt.result.wantErr {
			t.Errorf("getPublicKeyTypeUint8() error = %v, wantErr %v", err, tt.result.wantErr)
			return
		}
		assert.Equal(t, tt.result.pub_key_type, uPubKeyType)
	}
}

func TestGetPublicKeyTypeString(t *testing.T) {
	type args struct {
		pub_key_type uint8
	}
	type result struct {
		pub_key_type string
		wantErr      bool
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{
			"success",
			args{
				crypto.IDSECP256K1,
			},
			result{
				crypto.KindSECP256K1,
				false,
			},
		},
		{
			"err-pub_key_type-unknown",
			args{
				0,
			},
			result{
				"",
				true,
			},
		},
	}
	for _, tt := range tests {
		sPubKeyType, err := getPublicKeyTypeString(tt.args.pub_key_type)
		if (err != nil) != tt.result.wantErr {
			t.Errorf("getPublicKeyTypeString() error = %v, wantErr %v", err, tt.result.wantErr)
			return
		}
		assert.Equal(t, tt.result.pub_key_type, sPubKeyType)
	}
}
