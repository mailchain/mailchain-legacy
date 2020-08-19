package settings

import (
	"os"
	"testing"

	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"

	"github.com/golang/mock/gomock"
)

func TestCacheStoreProduce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Path    values.String
		Timeout values.String
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
					os.MkdirAll("./tmp/cache", os.ModePerm)
					m.EXPECT().Get().Return("./tmp/cache")
					return m
				}(),
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("10s")
					return m
				}(),
			},
			false,
			false,
		},
		{
			"invalid value for timeout",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					os.MkdirAll("./tmp/cache", os.ModePerm)
					m.EXPECT().Get().Return("./tmp/cache").Times(0)
					return m
				}(),
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("asdasda")
					return m
				}(),
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CacheStore{
				Path:    tt.fields.Path,
				Timeout: tt.fields.Timeout,
			}
			got, err := s.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("CacheStore.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("CacheStore.Produce() nil = %v, wantErr %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}
