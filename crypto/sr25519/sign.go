package sr25519

import (
	"crypto/rand"
	"errors"
	"io"

	"github.com/gtank/merlin"
	"github.com/gtank/ristretto255"
)

var substrateContext = []byte("substrate") //nolint gochecknoglobals

type signingContext struct {
	*merlin.Transcript
}

func newSigningContext(context, msg []byte) signingContext {
	transcript := merlin.NewTranscript("SigningContext")
	transcript.AppendMessage([]byte(""), context)
	transcript.AppendMessage([]byte("sign-bytes"), msg)

	return signingContext{transcript}
}

func (c *signingContext) challengeScalar(label []byte) *ristretto255.Scalar {
	b := c.ExtractBytes(label, 64)
	k := ristretto255.NewScalar()

	return k.FromUniformBytes(b)
}

func witness(nonce []byte) (*ristretto255.Scalar, error) {
	_ = nonce
	b := make([]byte, 64)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}

	return ristretto255.NewScalar().FromUniformBytes(b), nil
}

type signature struct {
	R *ristretto255.Element
	S *ristretto255.Scalar
}

func (s *signature) Encode() []byte { // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L91
	out := make([]byte, 64)
	copy(out[:32], s.R.Encode([]byte{}))
	copy(out[32:], s.S.Encode([]byte{}))
	out[63] |= 128

	return out
}

func (s *signature) Decode(sig []byte) error { // https://github.com/w3f/schnorrkel/blob/4112f6e8cb684a1cc6574f9097497e1e302ab9a8/src/sign.rs#L114
	if len(sig) != 64 {
		return errors.New("signature length must be 64")
	}

	s.R = ristretto255.NewElement()
	if err := s.R.Decode(sig[:32]); err != nil {
		return err
	}

	if sig[63]&128 == 0 {
		return errors.New("signature not marked as schnorrkel")
	}

	sig[63] &= 127
	s.S = ristretto255.NewScalar()

	return s.S.Decode(sig[32:])
}
