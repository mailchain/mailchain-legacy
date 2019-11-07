//nolint:dupl
package values

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func TestDefaultInt_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     int
		setting string
		store   Store
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			"use-value",
			fields{
				100,
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("setting-name").Return(true)
					m.EXPECT().GetInt("setting-name").Return(50)
					return m
				}(),
			},
			50,
		},
		{
			"use-default",
			fields{
				100,
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("setting-name").Return(false)
					return m
				}(),
			},
			100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultInt{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			if got := d.Get(); got != tt.want {
				t.Errorf("DefaultInt.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultInt_Set(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     int
		setting string
		store   Store
	}
	type args struct {
		v int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"set",
			fields{
				100,
				"setting-name",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().Set("setting-name", 100)
					return m
				}(),
			},
			args{
				100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultInt{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			d.Set(tt.args.v)
		})
	}
}

func TestNewDefaultInt(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		defVal  int
		store   Store
		setting string
	}
	tests := []struct {
		name string
		args args
		want Int
	}{
		{
			"success",
			args{
				100,
				valuestest.NewMockStore(nil),
				"setting",
			},
			DefaultInt{100, "setting", valuestest.NewMockStore(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultInt(tt.args.defVal, tt.args.store, tt.args.setting); !assert.Equal(tt.want, got) {
				t.Errorf("NewDefaultInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultInt_Attribute(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		def     int
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
				100,
				"test.setting.name1",
				func() Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("test.setting.name1").Return(true).AnyTimes()
					m.EXPECT().GetInt("test.setting.name1").Return(100).AnyTimes()
					return m
				}(),
			},
			output.Attribute{FullName: "name1", IsDefault: true, AdditionalComment: "", Value: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultInt{
				def:     tt.fields.def,
				setting: tt.fields.setting,
				store:   tt.fields.store,
			}
			if got := d.Attribute(); !assert.Equal(tt.want, got) {
				t.Errorf("DefaultInt.Attribute() = %v, want %v", got, tt.want)
			}
		})
	}
}
