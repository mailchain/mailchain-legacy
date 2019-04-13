package mailbox

import (
	"context"

	"github.com/gogo/protobuf/proto"
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
