package sr25519

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignVerify(t *testing.T) {
	tests := []struct {
		name         string
		signedBy     PrivateKey
		verifiedBy   PublicKey
		message      []byte
		wantErr      bool
		wantVerified bool
	}{
		{
			"bob-private-public-key",
			bobPrivateKey,
			bobPublicKey,
			[]byte("message"),
			false,
			true,
		},
		{
			"alice-private-public-key",
			alicePrivateKey,
			alicePublicKey,
			[]byte("egassem"),
			false,
			true,
		},
		{
			"alice-private-bob-public-key",
			alicePrivateKey,
			bobPublicKey,
			[]byte("egassem"),
			false,
			false,
		},
		{
			"bob-private-alice-public-key",
			bobPrivateKey,
			alicePublicKey,
			[]byte("message"),
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSig, err := tt.signedBy.Sign(tt.message)
			assert.Equal(t, tt.wantErr, err != nil)
			verified := tt.verifiedBy.Verify(tt.message, gotSig)
			assert.Equal(t, tt.wantVerified, verified)
		})
	}
}
