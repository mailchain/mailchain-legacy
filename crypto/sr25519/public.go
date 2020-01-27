package sr25519

import (
	"errors"

	"github.com/gtank/ristretto255"
	"github.com/mailchain/mailchain/crypto"
)

const (
	publicKeySize = 32
)

// PublicKey based on the sr25519 curve
type PublicKey struct {
	key []byte
}

// Bytes return Publickey Bytes
func (pk PublicKey) Bytes() []byte {
	return pk.key
}

// Kind returns the key type
func (pk PublicKey) Kind() string {
	return crypto.SR25519
}

// Verify uses the sr25519 signature algorithm to verify that the message was signed by
// this public key; it returns true if this key created the signature for the message,
// false otherwise
func (pk PublicKey) Verify(message, sig []byte) bool {
	signature := signature{}
	if err := signature.Decode(sig); err != nil {
		return false
	}

	context := newSigningContext(substrateContext, message)
	context.AppendMessage([]byte("proto-name"), []byte("Schnorr-sig"))
	context.AppendMessage([]byte("sign:pk"), pk.key)                      // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L212
	context.AppendMessage([]byte("sign:R"), signature.R.Encode([]byte{})) // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L213

	k := context.challengeScalar([]byte("sign:c")) // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L215
	// https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L216
	a := ristretto255.NewElement()
	if err := a.Decode(pk.key); err != nil {
		return false
	}

	Rp := ristretto255.NewElement()
	Rp = Rp.ScalarBaseMult(signature.S)
	ky := a.ScalarMult(k, a)
	Rp = Rp.Subtract(Rp, ky)

	return Rp.Equal(signature.R) == 1
}

// PublicKeyFromBytes - Convert byte array to PublicKey
func PublicKeyFromBytes(keyBytes []byte) (crypto.PublicKey, error) {
	switch len(keyBytes) {
	case publicKeySize:
		return &PublicKey{key: keyBytes}, nil
	default:
		return nil, errors.New("public key must be 32 bytes")
	}
}
