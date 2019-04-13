package mailbox

import "context"

// Sender signs a transaction the sends it
type Sender interface {
	Send(ctx context.Context, to []byte, from []byte, data []byte, signer Signer, opts SenderOpts) (err error)
}

// SenderOpts options for sending a message
type SenderOpts interface{}
