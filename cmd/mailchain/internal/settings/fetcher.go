package settings

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
)

func fetcher(s values.Store) *Fetcher {
	l := &Fetcher{
		Disabled: values.NewDefaultBool(true, s, "fetcher.disabled"),
	}

	return l
}

// Fetcher configuration element.
type Fetcher struct {
	Disabled values.Bool
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (l *Fetcher) Output() output.Element {
	return output.Element{
		FullName: "fetcher",
		Attributes: []output.Attribute{
			l.Disabled.Attribute(),
		},
	}
}
