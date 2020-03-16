package s3

import (
	"bytes"
	"context"
	"testing"

	"github.com/mailchain/mailchain/encoding"

	"github.com/mailchain/mailchain/stores/s3store"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewS3TransactionStore(t *testing.T) {
	type args struct {
		region string
		bucket string
		id     string
		secret string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			args{
				"region",
				"Bucket",
				"id",
				"secret",
			},
			false,
			false,
		},
		{
			"err-region",
			args{
				"",
				"Bucket",
				"id",
				"secret",
			},
			true,
			true,
		},
		{
			"err-Bucket",
			args{
				"region",
				"",
				"id",
				"secret",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewS3TransactionStore(tt.args.region, tt.args.bucket, tt.args.id, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil != tt.wantNil {
				t.Errorf("NewSent() got == nil = %v, wantNil %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestPutRawTransaction(t *testing.T) {
	type fields struct {
		uploader func(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
		bucket   string
	}
	type args struct {
		ctx            context.Context
		rawTransaction rawTransactionData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"success",
			fields{
				func(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(t, bytes.NewReader([]byte("{\"protocol\":\"tcp\",\"network\":\"mainnet\",\"hash\":\"0x636f6e74656e74732d68617368\",\"transaction\":{\"Txt\":\"data\"}}")), input.Body) {
						t.Errorf("body incorrect")
					}
					if !assert.Equal(t, aws.String("Bucket-id"), input.Bucket) {
						t.Errorf("Bucket incorrect")
					}
					if !assert.Equal(t, aws.String("636f6e74656e74732d68617368"), input.Key) {
						t.Errorf("Bucket incorrect")
					}

					return &s3manager.UploadOutput{Location: "https://Bucket-id/636f6e74656e74732d68617368"}, nil
				},
				"Bucket-id",
			},
			args{
				ctx: context.Background(),
				rawTransaction: rawTransactionData{
					Protocol: "tcp",
					Network:  "mainnet",
					Hash:     encoding.EncodeHexZeroX([]byte("contents-hash")),
					Tx: struct {
						Txt string
					}{
						Txt: "data",
					},
				},
			},
			false,
		},
		{
			"err-UploadProvider",
			fields{
				func(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(t, bytes.NewReader([]byte("{\"protocol\":\"tcp\",\"network\":\"mainnet\",\"hash\":\"0x636f6e74656e74732d68617368\",\"transaction\":{\"Txt\":\"data\"}}")), input.Body) {
						t.Errorf("body incorrect")
					}
					if !assert.Equal(t, aws.String("Bucket-id"), input.Bucket) {
						t.Errorf("Bucket incorrect")
					}
					if !assert.Equal(t, aws.String("636f6e74656e74732d68617368"), input.Key) {
						t.Errorf("Key incorrect")
					}

					return nil, errors.Errorf("failed to upload")
				},
				"Bucket-id",
			},
			args{
				ctx: context.Background(),
				rawTransaction: rawTransactionData{
					Protocol: "tcp",
					Network:  "mainnet",
					Hash:     encoding.EncodeHexZeroX([]byte("contents-hash")),
					Tx: struct {
						Txt string
					}{
						Txt: "data",
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := TransactionStore{
				uploader: &s3store.UploadProvider{
					Uploader: tt.fields.uploader,
					Bucket:   tt.fields.bucket,
				},
			}
			decodeHash, _ := encoding.DecodeHexZeroX(tt.args.rawTransaction.Hash)
			err := h.PutRawTransaction(tt.args.ctx, tt.args.rawTransaction.Protocol, tt.args.rawTransaction.Network,
				decodeHash, tt.args.rawTransaction.Tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sent.PutMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestKey(t *testing.T) {
	tests := []struct {
		name         string
		contentsHash []byte
		want         string
	}{
		{
			"success",
			[]byte("contents-hash"),
			encoding.EncodeHex([]byte("contents-hash")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := TransactionStore{}
			if got := h.Key(tt.contentsHash); got != tt.want {
				t.Errorf("Sent.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}
