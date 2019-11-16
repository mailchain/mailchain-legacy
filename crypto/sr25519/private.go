package sr25519

import (
	"github.com/pkg/errors"
	sr25519 "github.com/ChainSafe/go-schnorrkel"
)

const (
	chainCodeSize  = 32
	keyPairSize    = 96
	privateKeySize = 64
)

var SigningContext = []byte("substrate")

// Private Key sr25519
type PrivateKey struct {
	key *sr25519.SecretKey
}


// Sign uses the private key to sign the message using the sr25519 signature algorithm
func (k *PrivateKey) Sign(msg []byte) ([]byte, error) {
	if k.key == nil {
		return nil, errors.New("key is nil")
	}

	t := sr25519.NewSigningContext(SigningContext, msg)
	sig, err := k.key.Sign(t)
	if err != nil {
		return nil, err
	}

	enc := sig.Encode()
	return enc[:], nil
}

// Encode returns the 32-byte encoding of the private key
func (k *PrivateKey) Encode() []byte {
	if k.key == nil {
		return nil
	}

	enc := k.key.Encode()
	return enc[:]
}

// Decode decodes the input bytes into a private key and sets the receiver the decoded key
// Input must be 32 bytes, or else this function will error
func (k *PrivateKey) Decode(in []byte) error {
	if len(in) != privateKeySize {
		return errors.New("input to sr25519 private key decode is not 32 bytes")
	}
	b := [32]byte{}
	copy(b[:], in)
	k.key = &sr25519.SecretKey{}
	return k.key.Decode(b)
}

// PrivateKeyFromBytes get a private key from seed []byte
func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case privateKeySize:
		priv := &PrivateKey{key}
		b := [32]byte{}
		copy(b[:], privKey)

		private := priv.Decode(privKey)

		return &PrivateKey{
			key: &sr25519.SecretKey{ private }, 
			}, error

	default:
		return nil, errors.Errorf("sr25519: bad key length")			
	}
}

