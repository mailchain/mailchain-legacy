package schnorrkel

import "crypto/sha512"

// SecretKey consists of a secret scalar and a signing nonce
type SecretKey struct {
	seed  [32]byte
	key   [32]byte
	nonce [32]byte
}

func (sk *SecretKey) Key() []byte {
	return sk.key[:]
}

func (sk *SecretKey) Seed() []byte {
	return sk.seed[:]
}

func (sk *SecretKey) Nonce() []byte {
	return sk.nonce[:]
}

func NewSecretKeyED25519(seed [32]byte) SecretKey {
	key := [32]byte{}
	nonce := [32]byte{}
	h := sha512.Sum512(seed[:])

	copy(key[:], h[:32])

	key[0] &= 248
	key[31] &= 63
	key[31] |= 64
	t := divideScalarByCofactor(key[:])

	copy(key[:], t)
	copy(nonce[:], h[32:])

	return SecretKey{seed: seed, key: key, nonce: nonce}
}

// https://github.com/w3f/schnorrkel/blob/718678e51006d84c7d8e4b6cde758906172e74f8/src/scalars.rs#L18
func divideScalarByCofactor(s []byte) []byte {
	l := len(s) - 1 //nolint: gomnd length-1
	low := byte(0)

	for i := range s {
		r := s[l-i] & 0b00000111 //nolint: gomnd remainder
		s[l-i] >>= 3
		s[l-i] += low
		low = r << 5 //nolint: gomnd https://github.com/w3f/schnorrkel/blob/718678e51006d84c7d8e4b6cde758906172e74f8/src/scalars.rs#L34
	}

	return s
}
