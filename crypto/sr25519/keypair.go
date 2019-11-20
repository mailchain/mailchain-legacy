package sr25519

import (
	sr25519 "github.com/ChainSafe/go-schnorrkel"
)

// Keypair type have public and private key
type Keypair struct {
	public  *PublicKey
	private *PrivateKey
}

// NewKeypair returns a Sr25519 Keypair given a schnorrkel secret key
func NewKeypair(priv *sr25519.SecretKey) (*Keypair, error) {
	pub, err := priv.Public()
	if err != nil {
		return nil, err
	}

	return &Keypair{
		public:  &PublicKey{key: pub},
		private: &PrivateKey{key: priv},
	}, nil
}

// GenerateKeypair returns a new sr25519 keypair
func GenerateKeypair() (*Keypair, error) {
	priv, pub, err := sr25519.GenerateKeypair()
	if err != nil {
		return nil, err
	}

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
