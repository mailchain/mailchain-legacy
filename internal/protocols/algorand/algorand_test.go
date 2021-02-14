package algorand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworks(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			"success",
			[]string{"mainnet", "betanet", "testnet"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Networks(); !assert.Equal(t, tt.want, got) {
				t.Errorf("Networks() = %v, want %v", got, tt.want)
			}
		})
	}
}
