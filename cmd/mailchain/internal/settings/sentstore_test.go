package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/s3store"
	"github.com/stretchr/testify/assert"
)

func Test_sentStore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name     string
		args     args
		wantKind string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("sentstore.kind")
					return m
				}(),
			},
			"mailchain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sentStore(tt.args.s)
			assert.Equal(t, tt.wantKind, got.Kind.Get())
		})
	}
}

func TestSentStore_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Kind      values.String
		s3        *SentStoreS3
		mailchain *SentStoreMailchain
	}
	tests := []struct {
		name     string
		fields   fields
		wantType stores.Sent
		wantErr  bool
	}{
		{
			"s3",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("s3")
					return m
				}(),
				sentStoreS3(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("sentstore.s3.region").Return(true)
					m.EXPECT().GetString("sentstore.s3.region").Return("us-east-1")
					m.EXPECT().IsSet("sentstore.s3.accessKeyId").Return(false)
					m.EXPECT().IsSet("sentstore.s3.bucket").Return(true)
					m.EXPECT().GetString("sentstore.s3.bucket").Return("bucket-value")
					m.EXPECT().IsSet("sentstore.s3.secretAccessKey").Return(false)
					return m
				}()),
				nil,
			},
			&s3store.Sent{},
			false,
		},
		{
			"mailchain",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("mailchain")
					return m
				}(),
				nil,
				&SentStoreMailchain{},
			},
			&stores.SentStore{},
			false,
		},
		{
			"err",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("invalid").Times(2)
					return m
				}(),
				nil,
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := SentStore{
				Kind:      tt.fields.Kind,
				s3:        tt.fields.s3,
				mailchain: tt.fields.mailchain,
			}
			got, err := ss.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("SentStore.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.IsType(t, tt.wantType, got) {
				t.Errorf("SentStore.Produce() = %v, want %v", got, tt.wantType)
			}
		})
	}
}
