package settings

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/bdbstore"
	"github.com/mailchain/mailchain/stores/ldbstore"
	"github.com/stretchr/testify/assert"
)

func Test_mailboxState(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name     string
		args     args
		wantKind string
	}{
		{
			"check-defaults",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("mailboxState.kind").Return(false)
					return m
				}(),
			},
			"leveldb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mailboxState(tt.args.s)
			assert.Equal(tt.wantKind, got.Kind.Get())
		})
	}
}

func TestMailboxState_Produce(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Kind                 values.String
		mailboxStateLevelDB  MailboxStateLevelDB
		mailboxStateBadgerDB MailBoxStateBadgerDB
	}
	tests := []struct {
		name     string
		fields   fields
		wantType stores.State
		wantErr  bool
	}{
		{
			"success-leveldb",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("leveldb")
					return m
				}(),
				mailboxStateLevelDB(
					func() values.Store {
						m := valuestest.NewMockStore(mockCtrl)
						os.MkdirAll("./tmp/mailboxstate", os.ModePerm)
						m.EXPECT().IsSet("mailboxState.leveldb.path").Return(true)
						m.EXPECT().GetString("mailboxState.leveldb.path").Return("./tmp/mailboxstate")
						m.EXPECT().IsSet("mailboxState.leveldb.cache").Return(false)
						m.EXPECT().IsSet("mailboxState.leveldb.handles").Return(false)
						return m
					}(),
				),
				MailBoxStateBadgerDB{},
			},
			&ldbstore.Database{},
			false,
		},
		{
			"success-badgerdb",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("badgerdb")
					return m
				}(),
				MailboxStateLevelDB{},
				mailboxStateBadgerDB(
					func() values.Store {
						m := valuestest.NewMockStore(mockCtrl)
						os.MkdirAll("./tmp/mailboxstate", os.ModePerm)
						m.EXPECT().IsSet("mailboxState.badgerdb.path").Return(true)
						m.EXPECT().GetString("mailboxState.badgerdb.path").Return("./tmp/mailboxstate")
						return m
					}(),
				),
			},
			&bdbstore.Database{},
			false,
		},
		{
			"err-invalid",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("invalid").Times(2)
					return m
				}(),
				mailboxStateLevelDB(
					func() values.Store {
						m := valuestest.NewMockStore(mockCtrl)
						return m
					}(),
				),
				MailBoxStateBadgerDB{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MailboxState{
				Kind:                 tt.fields.Kind,
				mailboxStateLevelDB:  tt.fields.mailboxStateLevelDB,
				mailBoxStateBadgerDB: tt.fields.mailboxStateBadgerDB,
			}
			got, err := s.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("MailboxState.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.wantType, got) {
				t.Errorf("MailboxState.Produce() = %v, want %v", got, tt.wantType)
			}
		})
	}
}

func TestMailBoxStateBadgerDB_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Path values.String
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					os.MkdirAll("./tmp/badgerdb", os.ModePerm)
					m.EXPECT().Get().Return("./tmp/badgerdb")
					return m
				}(),
			},
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MailBoxStateBadgerDB{
				Path: tt.fields.Path,
			}

			got, err := s.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("MailboxStateBadgerDB.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("MailboxStateBadgerDB.Produce() nil = %v, wantErr %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestMailboxStateLevelDB_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Path    values.String
		Handles values.Int
		Cache   values.Int
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					os.MkdirAll("./tmp/leveldb", os.ModePerm)
					m.EXPECT().Get().Return("./tmp/leveldb")
					return m
				}(),
				func() values.Int {
					m := valuestest.NewMockInt(mockCtrl)
					m.EXPECT().Get().Return(0)
					return m
				}(),
				func() values.Int {
					m := valuestest.NewMockInt(mockCtrl)
					m.EXPECT().Get().Return(0)
					return m
				}(),
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MailboxStateLevelDB{
				Path:    tt.fields.Path,
				Handles: tt.fields.Handles,
				Cache:   tt.fields.Cache,
			}
			got, err := s.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("MailboxStateLevelDB.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("MailboxStateLevelDB.Produce() nil = %v, wantErr %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}
