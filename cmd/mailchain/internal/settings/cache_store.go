package settings

import (
	"time"

	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/cachestore"
)

func cacheStore(s values.Store) *CacheStore {
	k := &CacheStore{
		Path: values.NewDefaultString(defaults.MessageCachePath(), s, "cache.path"),
	}
	return k
}

type CacheStore struct {
	Path values.String
	//TODO: add duration
}

// Produce `stores.State` based on configuration settings.
func (s CacheStore) Produce() (stores.Cache, error) {
	return cachestore.NewCacheStore(1*time.Hour, s.Path.Get()), nil
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s CacheStore) Output() output.Element {
	return output.Element{
		FullName: "cache",
		Attributes: []output.Attribute{
			s.Path.Attribute(),
		},
	}
}
