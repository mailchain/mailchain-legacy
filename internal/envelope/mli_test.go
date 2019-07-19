package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMLIToAddress(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want map[uint64]string
	}{
		{
			"compatability",
			map[uint64]string{0x1: "https://mcx.mx"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MLIToAddress(); !assert.Equal(tt.want, got) {
				t.Errorf("MLIToAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
