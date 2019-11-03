//nolint:dupl
package values

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func TestDefaultStringSlice_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     []string
		setting string
		store   Store
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"use-value",
			fields{
				[]string{"def-1", "def-2"},
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("setting-name").Return(true)
					m.EXPECT().GetStringSlice("setting-name").Return([]string{"val-1", "val-2"})
					return m
				}(),
			},
			[]string{"val-1", "val-2"},
		},
		{
			"use-default",
			fields{
				[]string{"def-1", "def-2"},
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("setting-name").Return(false)
					return m
				}(),
			},
			[]string{"def-1", "def-2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultStringSlice{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			if got := d.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultStringSlice.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultStringSlice_Set(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     []string
		setting string
		store   Store
	}
	type args struct {
		v []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"set",
			fields{
				[]string{"val-1", "val-2"},
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().Set("setting-name", []string{"val-1", "val-2"})
					return m
				}(),
			},
			args{
				[]string{"val-1", "val-2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultStringSlice{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			d.Set(tt.args.v)
		})
	}
}

func TestNewDefaultStringSlice(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		defVal  []string
		store   Store
		setting string
	}
	tests := []struct {
		name string
		args args
		want StringSlice
	}{
		{
			"success",
			args{
				[]string{"a", "b"},
				valuestest.NewMockStore(nil),
				"setting",
			},
			DefaultStringSlice{[]string{"a", "b"}, "setting", valuestest.NewMockStore(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultStringSlice(tt.args.defVal, tt.args.store, tt.args.setting); !assert.Equal(tt.want, got) {
				t.Errorf("NewDefaultStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultStringSlice_Attribute(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     []string
		setting string
		store   Store
	}
	tests := []struct {
		name   string
		fields fields
		want   output.Attribute
	}{
		{
			"success",
			fields{
				[]string{"value1", "value2"},
				"test.setting.name1",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("test.setting.name1").Return(true).AnyTimes()
					m.EXPECT().GetStringSlice("test.setting.name1").Return([]string{"value1", "value2"}).AnyTimes()
					return m
				}(),
			},
			output.Attribute{FullName: "name1", IsDefault: true, AdditionalComment: "", Value: []string{"value1", "value2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultStringSlice{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			if got := d.Attribute(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultStringSlice.Attribute() = %v, want %v", got, tt.want)
			}
		})
	}
}
