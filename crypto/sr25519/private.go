package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/crypto"

	"github.com/pkg/errors"
)

const (
	chainCodeSize  = 32
	keyPairSize    = 96
	privateKeySize = 64
	seedSize       = 32
)

var SigningContext = []byte("substrate")

// Private Key sr25519
type PrivateKey struct {
	key *schnorrkel.SecretKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	b := pk.Encode()
	kb := make([]byte, len(b))
	copy(kb, b[:])
	return kb
}

// Kind is the type of private key.
func (pk PrivateKey) Kind() string {
	return crypto.SR25519
}

// input privatekey export PublickKey
func (pk PrivateKey) PublicKey() PublicKey {
	kp, err := NewKeypair(pk.key)
	if err != nil {
		panic(err)
	}
	return kp.Public()
}

// Sign uses the private key to sign the message using the sr25519 signature algorithm
func (k *PrivateKey) Sign(msg []byte) ([]byte, error) {
	if k.key == nil {
		return nil, errors.New("key is nil")
	}

	t := schnorrkel.NewSigningContext(SigningContext, msg)
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
	k.key = &schnorrkel.SecretKey{}
	return k.key.Decode(b)
}

func keyFromSeed(b []byte) (*schnorrkel.SecretKey, error) {
	kb := [32]byte{}
	copy(b, kb[:])

	priv, err := schnorrkel.NewMiniSecretKeyFromRaw(kb)
	if err != nil {
		return nil, err
	}

	return priv.ExpandUniform(), nil
}

// PrivateKeyFromBytes get a private key from seed []byte
func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case privateKeySize:
		privKey, err := keyFromSeed(privKey)
		if err != nil {
			return nil, err
		}
		return &PrivateKey{key: privKey}, nil
	case seedSize:
		privKey, err := keyFromSeed(privKey)
		if err != nil {
			return nil, err
		}
		return &PrivateKey{key: privKey}, nil
	case keyPairSize:
		privKey, err := keyFromSeed(privKey)
		if err != nil {
			return nil, err
		}
		pk, err := NewKeypair(privKey)
		if err != nil {
			return nil, err
		}
		return pk.private, nil
	default:
		return nil, errors.Errorf("bad key length")
	}
}
