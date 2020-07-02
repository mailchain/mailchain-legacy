package settings

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/stores/pinata"
)

func sentStorePinata(s values.Store) *SentStorePinata {
	return &SentStorePinata{
		APIKey:    values.NewDefaultString(defaults.Empty, s, "sentstore.pinata.api-key"),
		APISecret: values.NewDefaultString(defaults.Empty, s, "sentstore.pinata.api-secret"),
	}
}

// SentStorePinata configuration element.
type SentStorePinata struct {
	APIKey    values.String
	APISecret values.String
}

// Produce `s3store.Sent` based on configuration settings.
func (s SentStorePinata) Produce() (*pinata.Sent, error) {
	return pinata.NewSent(s.APIKey.Get(), s.APISecret.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s SentStorePinata) Output() output.Element {
	return output.Element{
		FullName: "sentstore.pinata",
		Attributes: []output.Attribute{
			s.APIKey.Attribute(),
			s.APISecret.Attribute(),
		},
	}
}
