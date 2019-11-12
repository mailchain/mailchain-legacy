package sr25519

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/developerfred/merlin"
	r255 "github.com/developerfred/ristretto255"
)

type MiniSecretKey struct {
	key *r255.Scalar
}

type PrivateKey struct {
	key   [32]byte
	nonce [32]byte
}

// NewRandomElement returns a random ristretto element
func NewRandomElement() (*r255.Element, error) {
	e := r255.NewElement()
	s := [64]byte{}
	_, err := rand.Read(s[:])
	if err != nil {
		return nil, err
	}

	return e.FromUniformBytes(s[:]), nil
}

// NewRandomScalar returns a random ristretto scalar
func NewRandomScalar() (*r255.Scalar, error) {
	s := [64]byte{}
	_, err := rand.Read(s[:])
	if err != nil {
		return nil, err
	}

	ss := r255.NewScalar()
	return ss.FromUniformBytes(s[:]), nil
}

// Public gets the public key corresponding to this secret key
func (s *PrivateKey) Public() (*PublicKey, error) {
	e := r255.NewElement()
	sc, err := ScalarFromBytes(s.key)
	if err != nil {
		return nil, err
	}
	return &PublicKey{key: e.ScalarBaseMult(sc)}, nil
}

// NewMiniSecretKey derives a mini secret key from a byte array
func NewMiniSecretKey(b [64]byte) *MiniSecretKey {
	s := r255.NewScalar()
	s.FromUniformBytes(b[:])
	return &MiniSecretKey{key: s}
}

// NewMiniSecretKeyFromRaw derives a mini secret key from little-endian encoded raw bytes.
func NewMiniSecretKeyFromRaw(b [32]byte) (*MiniSecretKey, error) {
	s := r255.NewScalar()
	err := s.Decode(b[:])
	if err != nil {
		return nil, err
	}

	s.Reduce()

	return &MiniSecretKey{key: s}, nil
}

// NewRandomMiniSecretKey generates a mini secret key from random
func NewRandomMiniSecretKey() (*MiniSecretKey, error) {
	s := [64]byte{}
	_, err := rand.Read(s[:])
	if err != nil {
		return nil, err
	}

	scpriv := r255.NewScalar()
	scpriv.FromUniformBytes(s[:])
	return &MiniSecretKey{key: scpriv}, nil
}

// ExpandUniform
func (s *MiniSecretKey) ExpandUniform() *PrivateKey {
	t := merlin.NewTranscript("ExpandSecretKeys")
	t.AppendMessage([]byte("mini"), s.key.Encode([]byte{}))
	scalarBytes := t.ExtractBytes([]byte("sk"), 64)
	key := r255.NewScalar()
	key.FromUniformBytes(scalarBytes[:])
	nonce := t.ExtractBytes([]byte("no"), 32)
	key32 := [32]byte{}
	copy(key32[:], key.Encode([]byte{}))
	nonce32 := [32]byte{}
	copy(nonce32[:], nonce)
	return &PrivateKey{
		key:   key32,
		nonce: nonce32,
	}
}

// PrivateKeyFromBytes get a private key from seed []byte
func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case ed25519.SeedSize:
		ExpandKey := privKey.ExpandUniform()
		return &PrivateKey{key: ExpandKey.key}, nil
	case ed25519.PrivateKeySize:
		secret := NewMiniSecretKey(privKey)
		pk := secret.ExpandUhiform()

		return &PrivateKey{key: pk, nonce: pk.nonce}
	default:
		return nil, erros.Errorf("sr25519: bad key length")
	}

}
