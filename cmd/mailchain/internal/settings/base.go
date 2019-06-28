package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/chains"
)

// TODO: maybe this is not a new?
func New(s values.Store) *Base {
	return &Base{
		Senders:          senders(s),
		Receivers:        receivers(s),
		PublicKeyFinders: publicKeyFinders(s),
		Protocols: map[string]*Protocol{
			chains.Ethereum: protocol(s, chains.Ethereum),
		},
		// other
		Keystore:     keystore(s),
		MailboxState: mailboxState(s),
		SentStore:    sentStore(s),
		Server:       server(s),
	}
}

type Base struct {
	Senders          *Senders
	Receivers        *Receivers
	PublicKeyFinders *PublicKeyFinders
	// protocol
	Protocols map[string]*Protocol
	// Ethereum  *Protocol
	// other
	Keystore     *Keystore
	MailboxState *MailboxState
	SentStore    *SentStore
	Server       *Server
}
