package aes256cbc

import (
	"encoding/hex"
	"testing"

	"github.com/mailchain/mailchain/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name     string
		original []byte
		expected []byte
		err      error
	}{
		{"no prefix:022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			testutil.MustHexDecodeString("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			testutil.MustHexDecodeString("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			nil,
		},
		{"with prefix:022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			testutil.MustHexDecodeString("042c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			testutil.MustHexDecodeString("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := compress(tc.original)
			assert.EqualValues(hex.EncodeToString(tc.expected), hex.EncodeToString(actual))
			assert.Equal(tc.err, err)
		})
	}
}

func TestDecompress(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name     string
		original []byte
		expected []byte
		err      error
	}{
		{"022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			testutil.MustHexDecodeString("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			testutil.MustHexDecodeString("042c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			nil,
		},
		{"03a34d6aef3eb42335fb3cacb59478c0b44c0bbeb8bb4ca427dbc7044157a5d24b",
			testutil.MustHexDecodeString("03a34d6aef3eb42335fb3cacb59478c0b44c0bbeb8bb4ca427dbc7044157a5d24b"),
			testutil.MustHexDecodeString("04a34d6aef3eb42335fb3cacb59478c0b44c0bbeb8bb4ca427dbc7044157a5d24b4adf14868d8449c9b3e50d3d6338f3e5a2d3445abe679cddbe75cb893475806f"),
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := decompress(tc.original)
			assert.EqualValues(hex.EncodeToString(tc.expected), hex.EncodeToString(actual))
			assert.Equal(tc.err, err)
		})
	}
}
