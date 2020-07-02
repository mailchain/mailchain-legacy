package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func Test_sentStorePinata(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name          string
		args          args
		wantAPIKey    string
		wantAPISecret string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("sentstore.pinata.api-key").Return(false)
					m.EXPECT().IsSet("sentstore.pinata.api-secret").Return(false)
					return m
				}(),
			},
			"",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sentStorePinata(tt.args.s)
			assert.Equal(t, tt.wantAPIKey, got.APIKey.Get())
			assert.Equal(t, tt.wantAPISecret, got.APISecret.Get())
		})
	}
}
