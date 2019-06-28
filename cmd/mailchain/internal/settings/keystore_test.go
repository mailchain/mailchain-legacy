package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	ks "github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/nacl"
	"github.com/stretchr/testify/assert"
)

func Test_keystore(t *testing.T) {
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
					m.EXPECT().IsSet("keystore.kind").Return(false)
					return m
				}(),
			},
			"nacl-filestore",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keystore(tt.args.s)
			assert.Equal(tt.wantKind, got.Kind.Get())
		})
	}
}

func TestKeystore_Produce(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Kind          values.String
		naclFileStore NACLFileStore
	}
	tests := []struct {
		name     string
		fields   fields
		wantType ks.Store
		wantErr  bool
	}{
		{
			"success",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("nacl-filestore")
					return m
				}(),
				naclFileStore(
					func() values.Store {
						m := valuestest.NewMockStore(mockCtrl)
						m.EXPECT().IsSet("keystore.nacl-filestore.path").Return(true)
						m.EXPECT().GetString("keystore.nacl-filestore.path").Return("./tmp")
						return m
					}(),
				),
			},
			&nacl.FileStore{},
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
				naclFileStore(
					func() values.Store {
						m := valuestest.NewMockStore(mockCtrl)
						return m
					}(),
				),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Keystore{
				Kind:          tt.fields.Kind,
				naclFileStore: tt.fields.naclFileStore,
			}
			got, err := s.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("Keystore.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.wantType, got) {
				t.Errorf("Keystore.Produce() = %v, want %v", got, tt.wantType)
			}
		})
	}
}

func TestNACLFileStore_Produce(t *testing.T) {
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
					m.EXPECT().Get().Return("./tmp")
					return m
				}(),
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NACLFileStore{
				Path: tt.fields.Path,
			}
			got, err := n.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("NACLFileStore.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("NACLFileStore.Produce() nil = %v, wantErr %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}
