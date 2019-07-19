package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKinds(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want map[byte]bool
	}{
		{
			"compatability",
			map[uint8]bool{0x1: true, 0x50: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Kinds(); !assert.Equal(tt.want, got) {
				t.Errorf("Kinds() = %v, want %v", got, tt.want)
			}
		})
	}
}
