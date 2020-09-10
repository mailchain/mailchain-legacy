package settings

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/bdbstore"
	"github.com/pkg/errors"
)

func mailboxState(s values.Store) *MailboxState {
	k := &MailboxState{
		Kind:                 values.NewDefaultString(defaults.MailboxStateKind, s, "mailboxState.kind"),
		mailBoxStateBadgerDB: mailboxStateBadgerDB(s),
	}
	return k
}

// MailboxState settings for mailbox state storage
type MailboxState struct {
	Kind                 values.String
	mailBoxStateBadgerDB MailBoxStateBadgerDB
}

// Produce `stores.State` based on configuration settings.
func (s MailboxState) Produce() (stores.State, error) {
	switch s.Kind.Get() {
	case StoreBadgerDB:
		return s.mailBoxStateBadgerDB.Produce()
	default:
		return nil, errors.Errorf("%q is an unsupported mailbox state", s.Kind.Get())
	}
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s MailboxState) Output() output.Element {
	return output.Element{
		FullName: "mailboxState",
		Elements: []output.Element{
			s.mailBoxStateBadgerDB.Output(),
		},
	}
}

func mailboxStateLevelDB(s values.Store) MailboxStateLevelDB {
	return MailboxStateLevelDB{
		Path:    values.NewDefaultString(defaults.MailboxStatePath(), s, "mailboxState.leveldb.path"),
		Handles: values.NewDefaultInt(0, s, "mailboxState.leveldb.handles"),
		Cache:   values.NewDefaultInt(256, s, "mailboxState.leveldb.cache"),
	}
}

// MailboxStateLevelDB settings
type MailboxStateLevelDB struct {
	Path    values.String
	Handles values.Int
	Cache   values.Int
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s MailboxStateLevelDB) Output() output.Element {
	return output.Element{
		FullName: "mailboxState.leveldb",
		Attributes: []output.Attribute{
			s.Path.Attribute(),
			s.Handles.Attribute(),
			s.Cache.Attribute(),
		},
	}
}

func mailboxStateBadgerDB(s values.Store) MailBoxStateBadgerDB {
	return MailBoxStateBadgerDB{
		Path: values.NewDefaultString(defaults.MailboxStatePath(), s, "mailboxState.badgerdb.path"),
	}
}

// MailboxStateBadgerDB settings
type MailBoxStateBadgerDB struct {
	Path values.String
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s MailBoxStateBadgerDB) Output() output.Element {
	return output.Element{
		FullName: "mailboxState.badgerdb",
		Attributes: []output.Attribute{
			s.Path.Attribute(),
		},
	}
}

// Produce a badgerdb database with settings applied.
func (s MailBoxStateBadgerDB) Produce() (*bdbstore.Database, error) {
	return bdbstore.New(s.Path.Get())
}
