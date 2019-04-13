package mailbox

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/mailchain/mailchain/internal/pkg/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/encoding"
	"github.com/pkg/errors"
)

// Sender signs a transaction the sends it
type Sender interface {
	Send(ctx context.Context, to []byte, from []byte, data []byte, signer Signer, opts SenderOpts) (err error)
}

// SenderOpts options for sending a message
type SenderOpts interface{}

func prefixedBytes(data proto.Message) ([]byte, error) {
	protoData, err := proto.Marshal(data)
	if err != nil {
		return nil, errors.WithMessage(err, "could not marshal data")
	}

	prefixedProto := make([]byte, len(protoData)+1)
	prefixedProto[0] = encoding.Protobuf
	copy(prefixedProto[1:], protoData)

	return prefixedProto, nil
}

// encryptLocation is encrypted with supplied public key and location string
func encryptLocation(pk keys.PublicKey, location string) ([]byte, error) {
	// TODO: encryptLocation hard coded to aes256cbc
	encryptedLocation, err := aes256cbc.Encrypt(pk, []byte(location))
	if err != nil {
		return nil, errors.WithMessage(err, "could not encrypt data")
	}
	return encryptedLocation, nil
}

// encryptMailMessage is encrypted with supplied public key and location string
func encryptMailMessage(pk keys.PublicKey, encodedMsg []byte) ([]byte, error) {
	// TODO: encryptMailMessage hard coded to aes256cbc
	encryptedData, err := aes256cbc.Encrypt(pk, encodedMsg)
	if err != nil {
		return nil, errors.WithMessage(err, "could not encrypt message")
	}

	return encryptedData, nil
}
