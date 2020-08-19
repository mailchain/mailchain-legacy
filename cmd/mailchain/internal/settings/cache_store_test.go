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
		Path values.String
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
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CacheStore{
				Path: tt.fields.Path,
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
