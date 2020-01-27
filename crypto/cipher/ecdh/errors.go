package ecdh

import "errors"

var (
	// ErrEphemeralGenerate generic error when failing to geneate ephemeral keys
	ErrEphemeralGenerate = errors.New("ecdh: ephemeral generation failed")
	// ErrSharedSecretGenerate generic error when failing to geneate a shared secret
	ErrSharedSecretGenerate = errors.New("ecdh: shared secret generation failed")
)
