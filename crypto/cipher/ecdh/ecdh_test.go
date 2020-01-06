package ecdh

import (
	"crypto/rand"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/stretchr/testify/assert"
)

func Test_SharedSecretEndToEnd(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		keyExchange         cipher.KeyExchange
		RecipientPrivateKey crypto.PrivateKey
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"secp256k1-random",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewSECP256K1(rand.Reader)
					return kx
				}(),
				func() crypto.PrivateKey {
					pk, err := secp256k1.GenerateKey(rand.Reader)
					if err != nil {
						assert.FailNow("secp256k1.GenerateKey error = %v", err)
					}
					return pk
				}(),
			},
		},
		{
			"secp256k1-charlotte",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewSECP256K1(rand.Reader)
					return kx
				}(),
				secp256k1test.CharlottePrivateKey,
			},
		},
		{
			"secp256k1-sofia",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewSECP256K1(rand.Reader)
					return kx
				}(),
				secp256k1test.SofiaPrivateKey,
			},
		},
		{
			"ed25519-random",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewED25519(rand.Reader)
					return kx
				}(),
				func() crypto.PrivateKey {
					pk, err := ed25519.GenerateKey(rand.Reader)
					if err != nil {
						assert.FailNow("ed25519.GenerateKey error = %v", err)
					}
					return pk
				}(),
			},
		},
		{
			"ed25519-sofia",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewED25519(rand.Reader)
					return kx
				}(),
				ed25519test.SofiaPrivateKey,
			},
		},
		{
			"ed25519-charlotte",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewED25519(rand.Reader)
					return kx
				}(),
				ed25519test.CharlottePrivateKey,
			},
		},
		{
			"sr25519-random",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewSR25519(rand.Reader)
					return kx
				}(),
				func() crypto.PrivateKey {
					pk, err := ed25519.GenerateKey(rand.Reader)
					if err != nil {
						assert.FailNow("sr25519.GenerateKey error = %v", err)
					}
					return pk
				}(),
			},
		},
		{
			"sr25519-sofia",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewED25519(rand.Reader)
					return kx
				}(),
				sr25519test.SofiaPrivateKey,
			},
		},
		{
			"sr25519-charlotte",
			args{
				func() cipher.KeyExchange {
					kx, _ := NewSR25519(rand.Reader)
					return kx
				}(),
				sr25519test.CharlottePrivateKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ephemeralPrivKey, err := tt.args.keyExchange.EphemeralKey()
			if err != nil {
				assert.Fail("EphemeralKey() error = %v", err)
				return
			}

			senderSharedSecret, err := tt.args.keyExchange.SharedSecret(ephemeralPrivKey, tt.args.RecipientPrivateKey.PublicKey())
			if err != nil {
				assert.Fail("SharedSecret() error = %v", err)
				return
			}
			recipientSharedSecret, err := tt.args.keyExchange.SharedSecret(tt.args.RecipientPrivateKey, ephemeralPrivKey.PublicKey())
			if err != nil {
				assert.Fail("SharedSecret() error = %v", err)
				return
			}
			controlPrivKey, err := tt.args.keyExchange.EphemeralKey()
			if err != nil {
				assert.Fail("EphemeralKey() error = %v", err)
				return
			}
			controlSenderSharedSecret, err := tt.args.keyExchange.SharedSecret(ephemeralPrivKey, controlPrivKey.PublicKey())
			if err != nil {
				assert.Fail("SharedSecret() error = %v", err)
				return
			}
			controlRecipientSharedSecret, err := tt.args.keyExchange.SharedSecret(controlPrivKey, ephemeralPrivKey.PublicKey())
			if err != nil {
				assert.Fail("SharedSecret() error = %v", err)
				return
			}
			assert.Equal(senderSharedSecret, recipientSharedSecret)
			assert.NotEqual(controlSenderSharedSecret, recipientSharedSecret)
			assert.NotEqual(controlRecipientSharedSecret, recipientSharedSecret)
		})
	}
}
