package stores

import (
	"fmt"

	"github.com/mailchain/mailchain/internal/pkg/crypto"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/pkg/errors"
)

// The Sender saves the message. This should not be used directly but as the first argument of storing.PutMessage.
type Sender interface {
	// PutMessage should write the message contents to the underlying storage service. Return the final location or any error.
	PutMessage(path string, msg []byte) (location string, err error)
}

// PutMessage does the pre work before saving the message as implemented by store.
func PutMessage(store Sender, messageID mail.ID, msg []byte) (location string, err error) {
	hash, err := crypto.CreateLocationHash(msg)
	if err != nil {
		return "", err
	}
	location, err = store.PutMessage(fmt.Sprintf("%s-%s", messageID.HexString(), hash.String()), msg)
	if err != nil {
		return "", errors.Wrap(err, "could not store message")
	}
	return location, nil
}
