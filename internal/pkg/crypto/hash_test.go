package crypto_test

import (
	"encoding/hex"
	"testing"

	"github.com/mailchain/mailchain/internal/pkg/crypto"
	"github.com/mailchain/mailchain/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCreateLocationHash(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name     string
		original []byte
		expected []byte
		err      error
	}{
		{"2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba",
			testutil.MustHexDecodeString("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			testutil.MustHexDecodeString("2204abd5fcd4"),
			nil,
		},
		{"022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			testutil.MustHexDecodeString("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			testutil.MustHexDecodeString("2204be6f4863"),
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := crypto.CreateLocationHash(tc.original)
			assert.EqualValues(hex.EncodeToString(tc.expected), hex.EncodeToString(actual))
			assert.Equal(tc.err, err)
		})
	}
}
