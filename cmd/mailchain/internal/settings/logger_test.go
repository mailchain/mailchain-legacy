package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func Test_logger(t *testing.T) {
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
			"check-defaults",
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
			if got := logger(tt.args.s); !assert.True(t, (got == nil) == tt.wantNil) {
				t.Errorf("logger() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}
