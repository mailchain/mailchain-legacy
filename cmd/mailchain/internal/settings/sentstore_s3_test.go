package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func Test_sentStoreS3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name                string
		args                args
		wantBucket          string
		wantRegion          string
		wantAccessKeyID     string
		wantSecretAccessKey string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("sentstore.s3.bucket").Return(false)
					m.EXPECT().IsSet("sentstore.s3.region").Return(false)
					m.EXPECT().IsSet("sentstore.s3.accessKeyId").Return(false)
					m.EXPECT().IsSet("sentstore.s3.secretAccessKey").Return(false)
					return m
				}(),
			},
			"",
			"",
			"",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sentStoreS3(tt.args.s)
			assert.Equal(t, tt.wantBucket, got.Bucket.Get())
			assert.Equal(t, tt.wantRegion, got.Region.Get())
			assert.Equal(t, tt.wantAccessKeyID, got.AccessKeyID.Get())
			assert.Equal(t, tt.wantSecretAccessKey, got.SecretAccessKey.Get())
		})
	}
}
