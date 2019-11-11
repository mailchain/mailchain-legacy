package sr25519

import (
	"crypto/rand"
	"crypto/sha512"

	"github.com/mailchain/mailchain/crypto"
	"github.com/developerfred/merlin"
	r255 "github.com/developerfred/ristretto255"
)

type MiniSecretKey struct {
	key *r255.Scalar
}

type PrivateKey struct {
	key   [32]byte
	nonce [32]byte
	Kind  [32]byte
}

// Kind is the type of private key.
func (pk PrivateKey) Kind() string {
	return crypto.ED25519
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

// Public gets the public key corresponding to this secret key
func (s *SecretKey) Public() (*PublicKey, error) {
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
func (s *MiniSecretKey) ExpandUniform() *SecretKey {
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
	return &SecretKey{
		key:   key32,
		nonce: nonce32,
	}
}

func (s *MiniSecretKey) ExpandEd25519() *SecretKey {
	h := sha512.Sum512(s.key.Encode([]byte{}))
	sk := &SecretKey{key: [32]byte{}, nonce: [32]byte{}}

	copy(sk.key[:], h[:32])
	sk.key[0] &= 248
	sk.key[31] &= 63
	sk.key[31] |= 64
	t := divideScalarByCofactor(sk.key[:])
	copy(sk.key[:], t)

	copy(sk.nonce[:], h[32:])

	return sk
}

// PrivateKeyFromBytes get a private key from seed []byte

func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case ed25519.SeedSize:
		return &PrivateKey{key: privKey.ExpandUniform()}, nil
	case ed25519.PrivateKeySize:
		secret := NewMiniSecretKey(privKey)
		pk := secret.ExpandUhiform()

		return &PrivateKey{key: pk}
	default:
		return nil, erros.Errorf("sr25519: bad key length")
	}

}
