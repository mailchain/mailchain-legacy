package sr25519

import (
	"fmt"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/crypto"
	"github.com/pkg/errors"
)

const (
	keyPairSize          = 96
	privateKeySize       = 64
	seedSize             = 32
	speedLength      int = 32
	privateKeyLength int = 32
	signatureLength  int = 64
)

// SigningContext sr25519
var SigningContext = []byte("substrate") //nolint gochecknoglobals

// PrivateKey sr25519
type PrivateKey struct {
	key *schnorrkel.SecretKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	b := pk.key.Encode()

	return b[:]
}

// Kind is the type of private key.
func (pk PrivateKey) Kind() string {
	return crypto.SR25519
}

// PublicKey return the crypto.PublicKey that is derived from the Privatekey
func (pk PrivateKey) PublicKey() crypto.PublicKey {
	secretKey := &(schnorrkel.SecretKey{})
	err := secretKey.Decode(pk.key.Encode())
	if err != nil {
		panic(fmt.Sprintf("Invalid private key: %v", err))
	}

	// sofia := encodingtest.MustDecodeHex("bfdef74a308d58ce9662781501e4ff7476f4bb86b3070306fa255726be8c867df7307cce368038060e53afc68ad19acef7600ec5eaaaeee4d33dee302661da3fad86c385610d5238fb91bdf021030e1b545efc2663d1dfb99296d6d20661204571b4c2a9f4994ba8a3e679433030989fa018a0c7044d674a2ec7b692c04bb7f7ac31401e5901d989d7e9c6fcbc70ec0a2162e06a30e7bbb423ea9c145f")
	// fmt.Println(sofia)
	pub, _ := secretKey.Public()

	return PublicKey{key: pub}
}

// Sign uses the PrivateKey to sign the message using the sr25519 signature algorithm
func (pk PrivateKey) Sign(message []byte) ([]byte, error) {
	if pk.key == nil {
		return nil, errors.New("key is nil")
	}

	t := schnorrkel.NewSigningContext(SigningContext, message)

	sig, err := pk.key.Sign(t)

	if err != nil {
		return nil, err
	}

	enc := sig.Encode()

	return enc[:], nil
}

// Encode returns the 32-byte encoding of the private key
func (pk PrivateKey) Encode() []byte {
	if pk.key == nil {
		return nil
	}

	enc := pk.key.Encode()

	return enc[:]
}

// Decode decodes the input bytes into a private key and sets the receiver the decoded key
// Input must be 32 bytes, or else this function will error
func (pk PrivateKey) Decode(in []byte) error {
	if len(in) != privateKeySize {
		return errors.New("input to sr25519 private key decode is not 32 bytes")
	}

	b := [32]byte{}
	copy(b[:], in)

	pk.key = &schnorrkel.SecretKey{}

	return pk.key.Decode(b)
}

func keyFromSeed(in []byte) (*schnorrkel.SecretKey, error) {
	if len(in) != seedSize {
		return nil, errors.New("input to sr25519 private key decode is not 32 bytes")
	}

	b := [32]byte{}
	copy(b[:], in)

	key := &schnorrkel.SecretKey{}
	err := key.Decode(b)

	return key, err
}

func keyFromBytes(in []byte) (*PrivateKey, error) {
	if len(in) != privateKeySize {
		return nil, errors.New("input to create sr25519 private key is no 64 bytes")
	}

	priv := new(PrivateKey)
	err := priv.Decode(in)

	return priv, err
}

// PrivateKeyFromBytes get a private key from seed []byte
func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case privateKeySize:
		privKey, err := keyFromBytes(privKey)
		if err != nil {
			return nil, err
		}

		return privKey, nil
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
		return nil, errors.Errorf("sr25519: bad key length")
	}
}
