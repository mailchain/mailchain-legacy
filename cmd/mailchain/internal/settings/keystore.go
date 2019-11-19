package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	ks "github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/nacl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus" //nolint: depguard
)

func keystore(s values.Store) *Keystore {
	k := &Keystore{
		Kind:          values.NewDefaultString(defaults.KeystoreKind, s, "keystore.kind"),
		naclFileStore: naclFileStore(s),
	}
	return k
}

// Keystore configuration element.
type Keystore struct {
	Kind          values.String
	naclFileStore NACLFileStore
}

// Produce `keystore.Store` based on configuration settings.
func (s Keystore) Produce() (ks.Store, error) {
	switch s.Kind.Get() {
	case StoreNACLFilestore:
		return s.naclFileStore.Produce()
	default:
		return nil, errors.Errorf("%q is an unsupported keystore", s.Kind.Get())
	}
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s Keystore) Output() output.Element {
	return output.Element{
		FullName: "keystore",
		Attributes: []output.Attribute{
			s.Kind.Attribute(),
		},
		Elements: []output.Element{
			s.naclFileStore.Output(),
		},
	}
}

func naclFileStore(s values.Store) NACLFileStore {
	return NACLFileStore{
		Path: values.NewDefaultString(defaults.KeystorePath(), s, "keystore.nacl-filestore.path"),
	}
}

// NACLFileStore configuration element.
type NACLFileStore struct {
	Path values.String
}

// Produce `nacl.FileStore` based on configuration settings.
func (n NACLFileStore) Produce() (*nacl.FileStore, error) {
	fs := nacl.NewFileStore(n.Path.Get(), logrus.StandardLogger().Writer())
	return &fs, nil
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (n NACLFileStore) Output() output.Element {
	return output.Element{
		FullName: "keystore.nacl-filestore",
		Attributes: []output.Attribute{
			n.Path.Attribute(),
		},
	}
}
