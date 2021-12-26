package sr25519

import (
	"io"

	"github.com/gtank/merlin"
	"github.com/gtank/ristretto255"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/internal/schnorrkel"
	"github.com/pkg/errors"
)

const (
	seedSize = 32
)

func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	seed := make([]byte, seedSize)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, err
	}

	return PrivateKeyFromBytes(seed)
}

// PrivateKey based on the sr25519 curve
type PrivateKey struct {
	secretKey schnorrkel.SecretKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	return pk.secretKey.Seed()
}

// PublicKey return the crypto.PublicKey that is derived from the Privatekey
func (pk PrivateKey) PublicKey() crypto.PublicKey {
	key := ristretto255.NewScalar()
	if err := key.Decode(pk.secretKey.Key()); err != nil {
		return nil
	}

	return &PublicKey{key: ristretto255.NewElement().ScalarBaseMult(key).Encode([]byte{})}
}

// Sign uses the PrivateKey to sign the message using the sr25519 signature algorithm
func (pk PrivateKey) Sign(message []byte) ([]byte, error) {
	context := newSigningContext(substrateContext, message)

	context.AppendMessage([]byte("proto-name"), []byte("Schnorr-sig")) // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L173
	context.AppendMessage([]byte("sign:pk"), pk.PublicKey().Bytes())   // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L174

	r, err := witness(pk.secretKey.Nonce()) // witness_scalar Not implemented https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L176
	if err != nil {
		return nil, err
	}

	R := ristretto255.NewElement().ScalarBaseMult(r)            // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L177
	context.AppendMessage([]byte("sign:R"), R.Encode([]byte{})) // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L179

	k := context.challengeScalar([]byte("sign:c")) // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L181

	pkScalar := ristretto255.NewScalar()
	if err := pkScalar.Decode(pk.secretKey.Key()); err != nil {
		return nil, err
	}

	s := pkScalar.Multiply(pkScalar, k).Add(pkScalar, r) // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L182
	sig := signature{R: R, S: s}

	return sig.Encode(), nil
}

// PrivateKeyFromBytes get a private key from seed []byte
func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case seedSize:
		seed := [seedSize]byte{}
		copy(seed[:], privKey)

		return &PrivateKey{secretKey: schnorrkel.NewSecretKeyED25519(seed)}, nil
	default:
		return nil, errors.Errorf("sr25519: bad key length")
	}
}

func ExchangeKeys(privKey *PrivateKey, pubKey *PublicKey, length int) ([]byte, error) {
	// https://github.com/w3f/schnorrkel/tree/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src
	transcript := merlin.NewTranscript("KEX")
	transcript.AppendMessage([]byte("ctx"), []byte{})
	privKey.secretKey.Key()

	a := ristretto255.NewElement()
	if err := a.Decode(pubKey.key); err != nil {
		return []byte{}, err
	}

	pkScalar := ristretto255.NewScalar()
	if err := pkScalar.Decode(privKey.secretKey.Key()); err != nil {
		return []byte{}, err
	}

	a.ScalarMult(pkScalar, a)
	transcript.AppendMessage([]byte{}, a.Encode([]byte{}))

	return transcript.ExtractBytes([]byte{}, length), nil
}
