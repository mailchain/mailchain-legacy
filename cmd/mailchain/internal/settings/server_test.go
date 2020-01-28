package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func Test_cors(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name               string
		args               args
		wantDisabled       bool
		wantAllowedOrigins []string
	}{
		{
			"check-defaults",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("server.cors.disabled").Return(false)
					m.EXPECT().IsSet("server.cors.allowedOrigins").Return(false)
					return m
				}(),
			},
			false,
			[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cors(tt.args.s)
			assert.Equal(tt.wantDisabled, got.Disabled.Get())
			assert.Equal(tt.wantAllowedOrigins, got.AllowedOrigins.Get())
		})
	}
}

func Test_server(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name     string
		args     args
		wantPort int
	}{
		{
			"check-defaults",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("server.port").Return(false)
					return m
				}(),
			},
			8080,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server(tt.args.s)
			assert.Equal(tt.wantPort, got.Port.Get())
		})
	}
}
