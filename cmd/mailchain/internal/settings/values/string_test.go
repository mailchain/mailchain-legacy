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

func TestDefaultString_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     string
		setting string
		store   Store
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"use-value",
			fields{
				"def-1",
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("setting-name").Return(true)
					m.EXPECT().GetString("setting-name").Return("val-1")
					return m
				}(),
			},
			"val-1",
		},
		{
			"use-default",
			fields{
				"def-1",
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("setting-name").Return(false)
					return m
				}(),
			},
			"def-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultString{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			if got := d.Get(); got != tt.want {
				t.Errorf("DefaultString.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultString_Set(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     string
		setting string
		store   Store
	}
	type args struct {
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"set",
			fields{
				"val-1",
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().Set("setting-name", "val-1")
					return m
				}(),
			},
			args{
				"val-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultString{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			d.Set(tt.args.v)
		})
	}
}

func TestNewDefaultString(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		defVal  string
		store   Store
		setting string
	}
	tests := []struct {
		name string
		args args
		want String
	}{
		{
			"success",
			args{
				"value",
				valuestest.NewMockStore(nil),
				"setting",
			},
			DefaultString{"value", "setting", valuestest.NewMockStore(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultString(tt.args.defVal, tt.args.store, tt.args.setting); !assert.Equal(tt.want, got) {
				t.Errorf("NewDefaultString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultString_Attribute(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     string
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
				"value1",
				"test.setting.name1",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("test.setting.name1").Return(true).AnyTimes()
					m.EXPECT().GetString("test.setting.name1").Return("value1").AnyTimes()
					return m
				}(),
			},
			output.Attribute{FullName: "name1", IsDefault: true, AdditionalComment: "", Value: "value1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultString{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			if got := d.Attribute(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultString.Attribute() = %v, want %v", got, tt.want)
			}
		})
	}
}
