package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
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

// SentStore configuration element
type SentStore struct {
	Kind      values.String
	s3        *SentStoreS3
	mailchain *SentStoreMailchain
}

// Produce `stores.Sent` based on configuration settings.
func (ss SentStore) Produce() (stores.Sent, error) {
	switch ss.Kind.Get() {
	case StoreS3:
		return ss.s3.Produce()
	case defaults.Mailchain:
		return ss.mailchain.Produce()
	default:
		return nil, errors.Errorf("%q is an unsupported sent store", ss.Kind.Get())
	}
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (ss SentStore) Output() output.Element {
	return output.Element{
		FullName:   "sentstore",
		Attributes: []output.Attribute{ss.Kind.Attribute()},
		Elements: []output.Element{
			ss.s3.Output(),
		},
	}
}
