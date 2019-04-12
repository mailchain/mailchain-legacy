package aes256cbc

import "github.com/pkg/errors"

type encryptedData struct {
	InitializationVector      []byte `json:"iv"`
	EphemeralPublicKey        []byte `json:"ephemPublicKey"`
	Ciphertext                []byte `json:"ciphertext"`
	MessageAuthenticationCode []byte `json:"mac"`
}

func (e *encryptedData) verify() error {
	if len(e.InitializationVector) != 16 {
		return errors.Errorf("`InitializationVector` must be 16")
	}
	if len(e.EphemeralPublicKey) != 65 {
		return errors.Errorf("`EphemeralPublicKey` must be 65")
	}
	if len(e.MessageAuthenticationCode) != 32 {
		return errors.Errorf("`MessageAuthenticationCode` must be 16")
	}
	if len(e.Ciphertext) == 0 {
		return errors.Errorf("`Ciphertext` must not be empty")
	}

	return nil
}
