package pq

import (
	"testing"

	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

func TestGetProtocolNetworkUint8(t *testing.T) {
	assert := assert.New(t)
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
			"protocol unknown",
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
			"network unknown",
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
		assert.Equal(tt.result.protocol, uProtocol)
		assert.Equal(tt.result.network, uNetwork)
	}
}
