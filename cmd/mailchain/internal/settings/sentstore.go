package settings

import (
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

func sentStore(s values.Store) *SentStore {
	ss := &SentStore{
		Kind:      values.NewDefaultString(defaults.SentStoreKind, s, "sentstore.kind"),
		s3:        sentStoreS3(s),
		mailchain: &SentStoreMailchain{},
	}
	return ss
}

type SentStore struct {
	Kind      values.String
	s3        *SentStoreS3
	mailchain *SentStoreMailchain
}

func (ss SentStore) Produce() (stores.Sent, error) {
	switch ss.Kind.Get() {
	case mailchain.StoreS3:
		return ss.s3.Produce()
	case mailchain.Mailchain:
		return ss.mailchain.Produce()
	default:
		return nil, errors.Errorf("%q is an unsupported sent store", ss.Kind.Get())
	}
}
