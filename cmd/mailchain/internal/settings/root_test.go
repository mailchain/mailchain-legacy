package settings

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update .golden files")

func TestFromStore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					return m
				}(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromStore(tt.args.s)
			if (got == nil) != tt.wantNil {
				t.Errorf("FromStore() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestRoot_ToYaml(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s               values.Store
		tabsize         int
		commentDefaults bool
		excludeDefaults bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"comment-defaults",
			args{
				func() values.Store {
					v := viper.New()
					v.Set("keystore.nacl-filestore.path", "/home/user/.mailchain/.keystore")
					v.Set("mailboxState.leveldb.path", "/home/user/.mailchain/.mailbox")
					return v
				}(),
				2,
				true,
				false,
			},
		},
		{
			"exclude-defaults",
			args{
				func() values.Store {
					v := viper.New()
					v.Set("server.port", 12345)
					return v
				}(),
				2,
				false,
				true,
			},
		},
		{
			"include-defaults",
			args{
				func() values.Store {
					v := viper.New()
					v.Set("keystore.nacl-filestore.path", "/home/user/.mailchain/.keystore")
					v.Set("mailboxState.leveldb.path", "/home/user/.mailchain/.mailbox")
					v.Set("server.port", 12345)
					return v
				}(),
				2,
				false,
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromStore(tt.args.s)
			if got == nil {
				t.Errorf("FromStore() nil = %v, wantNil false", got == nil)
				return
			}
			out := &bytes.Buffer{}
			got.ToYaml(out, tt.args.tabsize, tt.args.commentDefaults, tt.args.excludeDefaults)

			golden := filepath.Join("./testdata/", t.Name()+"."+tt.name+".golden.yaml")
			want, _ := ioutil.ReadFile(golden)

			if *update {
				err := ioutil.WriteFile(golden, out.Bytes(), 0644)
				assert.NoError(err)
			}

			assert.EqualValues(string(want), out.String())
		})
	}
}
