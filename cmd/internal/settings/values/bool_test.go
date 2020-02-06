package values

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func TestDefaultBool_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     bool
		setting string
		store   Store
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"use-value",
			fields{
				false,
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("setting-name").Return(true)
					m.EXPECT().GetBool("setting-name").Return(true)
					return m
				}(),
			},
			true,
		},
		{
			"use-default",
			fields{
				true,
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("setting-name").Return(false)
					return m
				}(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultBool{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			if got := d.Get(); got != tt.want {
				t.Errorf("DefaultBool.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBool_Set(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     bool
		setting string
		store   Store
	}
	type args struct {
		v bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"set",
			fields{
				true,
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().Set("setting-name", true)
					return m
				}(),
			},
			args{
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultBool{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			d.Set(tt.args.v)
		})
	}
}

func TestNewDefaultBool(t *testing.T) {
	type args struct {
		defVal  bool
		store   Store
		setting string
	}
	tests := []struct {
		name string
		args args
		want Bool
	}{
		{
			"success",
			args{
				true,
				valuestest.NewMockStore(nil),
				"setting",
			},
			DefaultBool{true, "setting", valuestest.NewMockStore(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultBool(tt.args.defVal, tt.args.store, tt.args.setting); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewDefaultBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBool_Attribute(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     bool
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
				false,
				"test.setting.name1",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("test.setting.name1").Return(true).AnyTimes()
					m.EXPECT().GetBool("test.setting.name1").Return(true).AnyTimes()
					return m
				}(),
			},
			output.Attribute{FullName: "name1", IsDefault: false, AdditionalComment: "", Value: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultBool{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			if got := d.Attribute(); !assert.Equal(t, tt.want, got) {
				t.Errorf("DefaultBool.Attribute() = %v, want %v", got, tt.want)
			}
		})
	}
}
