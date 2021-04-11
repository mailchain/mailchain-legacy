package sr25519

import (
	"errors"

	"github.com/gtank/ristretto255"
	"github.com/mailchain/mailchain/crypto"
)

const (
	balKeySize = 32
)

// PublicKey based on the sr25519 curve
type Bal struct {
	key []byte
}

// Bytes return Publickey Bytes
func (b Bal) Bytes() []byte {
	return b.key
}

// Kind returns the key type
func (b Bal) Kind() string {
	return crypto.KindSR25519
}

// Kind returns the key type
func (b Bal) Balance() string {
	return crypto.KindSR25519
}

// Kind returns the key type
func (b Bal) Unit() string {
	return crypto.KindSR25519
}

// Verify uses the sr25519 signature algorithm to verify that the message was signed by
// this public key; it returns true if this key created the signature for the message,
// false otherwise
func (b Bal) Verify(message, sig []byte) bool {
	signature := signature{}
	if err := signature.Decode(sig); err != nil {
		return false
	}

	context := newSigningContext(substrateContext, message)
	context.AppendMessage([]byte("proto-name"), []byte("Schnorr-sig"))
	context.AppendMessage([]byte("sign:pk"), b.key)                       // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L212
	context.AppendMessage([]byte("sign:R"), signature.R.Encode([]byte{})) // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L213

	k := context.challengeScalar([]byte("sign:c")) // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L215
	// https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L216
	a := ristretto255.NewElement()
	if err := a.Decode(b.key); err != nil {
		return false
	}

	Rp := ristretto255.NewElement()
	Rp = Rp.ScalarBaseMult(signature.S)
	ky := a.ScalarMult(k, a)
	Rp = Rp.Subtract(Rp, ky)

	return Rp.Equal(signature.R) == 1
}

// PublicKeyFromBytes - Convert byte array to PublicKey
func BalanceFromBytes(keyBytes []byte) (crypto.Balance, error) {
	switch len(keyBytes) {
	case balKeySize:
		return &Bal{key: keyBytes}, nil
	default:
		return nil, errors.New("balance must be 32 bytes")
	}
}
