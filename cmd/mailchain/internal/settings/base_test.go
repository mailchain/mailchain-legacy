package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
)

func TestNew(t *testing.T) {
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
			got := New(tt.args.s)
			if (got == nil) != tt.wantNil {
				t.Errorf("New() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}
