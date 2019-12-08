package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/pkg/errors"
)

// Keypair type have public and private key
type Keypair struct {
	public  *PublicKey
	private *PrivateKey
}

// NewKeypair returns a Sr25519 Keypair given a schnorrkel secret key
func NewKeypair(priv *schnorrkel.SecretKey) (*Keypair, error) {
	pub, err := priv.Public()
	if err != nil {
		return nil, err
	}

	return &Keypair{
		public:  &PublicKey{key: pub},
		private: &PrivateKey{key: priv},
	}, nil
}

// NewPublicKey creates a new public key using the input bytes
func NewPublicKey(in []byte) (*PublicKey, error) {
	if len(in) != publicKeyLength {
		return nil, errors.New("cannot create public key: input is not 32 bytes")
	}

	buf := [publicKeyLength]byte{}
	copy(buf[:], in)

	return &PublicKey{key: schnorrkel.NewPublicKey(buf)}, nil
}

// NewPrivateKey creates a new private key using the input bytes
func NewPrivateKey(in []byte) (*PrivateKey, error) {
	if len(in) != privateKeyLength {
		return nil, errors.New("input to create sr25519 private key is not 32 bytes")
	}

	priv := new(PrivateKey)
	err := priv.Decode(in)

	return priv, err
}

// GenerateKeypair returns a new sr25519 keypair
func GenerateKeypair() (*Keypair, error) {
	priv, pub, err := schnorrkel.GenerateKeypair()
	if err != nil {
		return nil, err
	}

	return &Keypair{
		public:  &PublicKey{key: pub},
		private: &PrivateKey{key: priv},
	}, nil
}

// NewKeypairFromSeed returns a new Keypair given a seed
func NewKeypairFromSeed(seed []byte) (*Keypair, error) {
	buf := [32]byte{}

	msc, err := schnorrkel.NewMiniSecretKeyFromRaw(buf)
	if err != nil {
		return nil, err
	}

	priv := msc.ExpandEd25519()
	pub := msc.Public()

	return &Keypair{
		public:  &PublicKey{key: pub},
		private: &PrivateKey{key: priv},
	}, nil
}

// Sign uses the keypair to sign the message using the sr25519 signature algorithm
func (kp *Keypair) Sign(msg []byte) ([]byte, error) {
	return kp.private.Sign(msg)
}

// Public returns the public key corresponding to this keypair
func (kp *Keypair) Public() *PublicKey {
	return kp.public
}

// Private returns the private key corresponding to this keypair
func (kp *Keypair) Private() *PrivateKey {
	return kp.private
}
