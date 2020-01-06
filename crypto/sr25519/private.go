package sr25519

import (
	"crypto/ed25519"
	"io"

	"github.com/developerfred/go-schnorrkel"
	"github.com/mailchain/mailchain/crypto"
	"github.com/pkg/errors"
)

const (
	privateKeySize   = 64
	seedSize         = 32
	privateKeyLength = 32
)

// SigningContext sr25519
var SigningContext = []byte("substrate") //nolint gochecknoglobals

func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	_, pPrivKey, err := ed25519.GenerateKey(rand)
	if err != nil {
		return nil, err
	}

	return PrivateKeyFromBytes(pPrivKey)

}

// PrivateKey sr25519
type PrivateKey struct {
	key []byte
}

func (pk PrivateKey) generate() (*schnorrkel.SecretKey, error) {
	if pk.key == nil {
		return &schnorrkel.SecretKey{}, errors.New("invalid key")
	}

	b := [32]byte{}
	copy(b[:], pk.key)

	msc, err := schnorrkel.NewMiniSecretKeyFromRaw(b)
	if err != nil {
		return nil, err
	}

	return msc.ExpandEd25519(), nil
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	return pk.key
}

// Kind is the type of private key.
func (pk PrivateKey) Kind() string {
	return crypto.SR25519
}

// PublicKey return the crypto.PublicKey that is derived from the Privatekey
func (pk PrivateKey) PublicKey() crypto.PublicKey {
	if pk.key == nil {
		return PublicKey{}
	}

	msc, _ := pk.generate()

	public, _ := msc.Public()
	pb := public.Encode()

	return PublicKey{key: pb[:]}
}

// Sign uses the PrivateKey to sign the message using the sr25519 signature algorithm
func (pk PrivateKey) Sign(message []byte) (signature []byte, err error) {
	if pk.key == nil {
		return nil, errors.New("cannot create private key: input is not 32 bytes")
	}

	priv, _ := pk.generate()

	signingContext := schnorrkel.NewSigningContext(SigningContext, message)

	sig, err := priv.Sign(signingContext)
	if err != nil {
		return []byte{}, err
	}

	enc := sig.Encode()

	return enc[:], nil
}

func keyFromSeed(in []byte) (*schnorrkel.SecretKey, error) {
	if len(in) != seedSize {
		return nil, errors.New("input to sr25519 private key decode is not 32 bytes")
	}

	b := [privateKeyLength]byte{}
	copy(b[:], in)

	key := schnorrkel.SecretKey{}
	err := key.Decode(b)

	return &key, err
}

// PrivateKeyFromBytes get a private key from seed []byte
func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case privateKeySize:
		return &PrivateKey{key: privKey}, nil
	case seedSize:
		k, _ := keyFromSeed(privKey)
		b := k.Encode()

		return &PrivateKey{key: b[:]}, nil
	default:
		return nil, errors.Errorf("sr25519: bad key length")
	}
}
