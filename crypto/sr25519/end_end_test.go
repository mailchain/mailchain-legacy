package sr25519

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignVerify(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		signedBy     PrivateKey
		verifiedBy   PublicKey
		message      []byte
		wantErr      bool
		wantVerified bool
	}{
		{
			"charlotte-private-public-key",
			charlottePrivateKey,
			charlottePublicKey,
			[]byte("message"),
			false,
			true,
		},
		{
			"sofia-private-public-key",
			sofiaPrivateKey,
			sofiaPublicKey,
			[]byte("egassem"),
			false,
			true,
		},
		{
			"sofia-private-charlotte-public-key",
			sofiaPrivateKey,
			charlottePublicKey,
			[]byte("egassem"),
			false,
			false,
		},
		{
			"charlotte-private-sofia-public-key",
			charlottePrivateKey,
			sofiaPublicKey,
			[]byte("message"),
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSig, err := tt.signedBy.Sign(tt.message)
			assert.Equal(tt.wantErr, err != nil)
			verified := tt.verifiedBy.Verify(tt.message, gotSig)
			assert.Equal(tt.wantVerified, verified)
		})
	}
}
