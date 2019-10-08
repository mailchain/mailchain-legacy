package settings

import (
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/ldbstore"
	"github.com/pkg/errors"
)

func mailboxState(s values.Store) *MailboxState {
	k := &MailboxState{
		Kind:                values.NewDefaultString(defaults.MailboxStateKind, s, "mailboxState.kind"),
		mailboxStateLevelDB: mailboxStateLevelDB(s),
	}
	return k
}

type MailboxState struct {
	Kind                values.String
	mailboxStateLevelDB MailboxStateLevelDB
}

func (s MailboxState) Produce() (stores.State, error) {
	switch s.Kind.Get() {
	case mailchain.StoreLevelDB:
		return s.mailboxStateLevelDB.Produce()
	default:
		return nil, errors.Errorf("%q is an unsupported mailbox state", s.Kind.Get())
	}
}

func (s MailboxState) Output() output.Element {
	return output.Element{
		FullName: "mailboxState",
		Elements: []output.Element{
			s.mailboxStateLevelDB.Output(),
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

type MailboxStateLevelDB struct {
	Path    values.String
	Handles values.Int
	Cache   values.Int
}

func (s MailboxStateLevelDB) Produce() (*ldbstore.Database, error) {
	return ldbstore.New(s.Path.Get(), s.Cache.Get(), s.Handles.Get())
}

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
