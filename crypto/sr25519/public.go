package sr25519

import (
	"crypto/ed25519"
	"crypto/sha512"
	" github.com/mailchain/mailchain/crypto/sr25519"

	r255 "github.com/developerfred/ristretto255"
	"github.com/pkg/errors"
)

// PublicKey is a member
type PublicKey struct {
	key *r255.Element
}
type MiniSecretKey struct {
	key *r255.Scalar
}

func divideScalarByCofactor(s []byte) []byte {
	l := len(s) - 1
	low := byte(0)
	for i := range s {
		r := s[l-i] & 0b00000111 // remainder
		s[l-i] >>= 3
		s[l-i] += low
		low = r << 5
	}

	return s
}

// ExpandEd25519 expands a mini PrivateKey key into a private key
func (s *MiniSecretKey) ExpandEd25519() *sr25519.PrivateKey {
	h := sha512.Sum512(s.key.Encode([]byte{}))
	sk := &sr25519.PrivateKey{key: [32]byte{}, nonce: [32]byte{}}

	copy(sk.key[:], h[:32])
	sk.key[0] &= 248
	sk.key[31] &= 63
	sk.key[31] |= 64
	t := divideScalarByCofactor(sk.key[:])
	copy(sk.key[:], t)

	copy(sk.nonce[:], h[32:])

	return sk
}

// Public gets the public key corresponding to this mini private key
func (s *MiniSecretKey) Public() *PublicKey {
	e := r255.NewElement()
	sk := s.ExpandEd25519()
	skey, err := ScalarFromBytes(sk.key)
	if err != nil {
		return nil
	}
	return &PublicKey{key: e.ScalarBaseMult(skey)}
}

// ScalarFromBytes returns a ristretto scalar from the input bytes
// performs input mod l where l is the group order
func ScalarFromBytes(b [32]byte) (*r255.Scalar, error) {
	s := r255.NewScalar()
	err := s.Decode(b[:])
	if err != nil {
		return nil, err
	}

	s.Reduce()
	return s, nil
}

// Compress returns the encoding of the point underlying the public key
func (p *PublicKey) Compress() [32]byte {
	b := p.key.Encode([]byte{})
	enc := [32]byte{}
	copy(enc[:], b)
	return enc
}

// Convert this public key to a byte array.
func PublicKeyFromBytes(keyBytes []byte) (*PublicKey, error) {
	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, errors.Errorf("public key must be 32 bytes")
	}
	pub := r255.NewElement()
	pub.FromUniformBytes(keyBytes[:])

	return &PublicKey{key: pub}, nil
}
