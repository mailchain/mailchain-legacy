package sr25519

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/sr25519"
	"github.com/mailchain/mailchain/crypto"

	"github.com/gtank/merlin"
	r255 "github.com/gtank/ristreto255"
	
)

// PublicKey is a member
type PublicKey struct {
	key *r255.Element
}

// Public gets the public key corresponding to this mini secret key
func (s *MiniSecretKey) Public() *PublicKey {
	e := r255.NewElement()
	sk := s.ExpandEd25519()
	skey, err := ScalarFromBytes(sk.key)
	if err != nil {
		return nil
	}
	return &PublicKey{key: e.ScalarBaseMult(skey)}
}

// Public gets the public key corresponding to this secret key
func (s *SecretKey) Public() (*PublicKey, error) {
	e := r255.NewElement()
	sc, err := ScalarFromBytes(s.key)
	if err != nil {
		return nil, err
	}
	return &PublicKey{key: e.ScalarBaseMult(sc)}, nil
}

// Compress returns the encoding of the point underlying the public key
func (p *PublicKey) Compress() [32]byte {
	b := p.key.Encode([]byte{})
	enc := [32]byte{}
	copy(enc[:], b)
	return enc
}