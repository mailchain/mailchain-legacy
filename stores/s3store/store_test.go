package s3store

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mailchain/mailchain/encoding"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewUploader(t *testing.T) {
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
			got, err := NewUploader(tt.args.region, tt.args.bucket, tt.args.id, tt.args.secret)
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

func TestUpload(t *testing.T) {

	type fields struct {
		uploader func(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
		bucket   string
	}
	type args struct {
		contentsHash []byte
		msg          []byte
		headers      map[string]*string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantAddress string
		wantErr     bool
	}{
		{
			"success-no-headers",
			fields{
				func(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(t, bytes.NewReader([]byte("test-data")), input.Body) {
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
				[]byte("contents-hash"),
				[]byte("test-data"),
				nil,
			},
			"https://Bucket-id/636f6e74656e74732d68617368",
			false,
		},
		{
			"success-has-headers",
			fields{
				func(ctx aws.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(t, bytes.NewReader([]byte("test-data")), input.Body) {
						t.Errorf("body incorrect")
					}
					if !assert.Equal(t, aws.String("Bucket-id"), input.Bucket) {
						t.Errorf("Bucket incorrect")
					}
					if !assert.Equal(t, aws.String("636f6e74656e74732d68617368"), input.Key) {
						t.Errorf("Key incorrect")
					}

					return &s3manager.UploadOutput{Location: "https://Bucket-id/636f6e74656e74732d68617368"}, nil
				},
				"Bucket-id",
			},
			args{
				[]byte("contents-hash"),
				[]byte("test-data"),
				map[string]*string{
					"Key-1": aws.String("value-1"),
				},
			},
			"https://Bucket-id/636f6e74656e74732d68617368",
			false,
		},
		{
			"err-UploadProvider",
			fields{
				func(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(t, bytes.NewReader([]byte("test-data")), input.Body) {
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
				[]byte("contents-hash"),
				[]byte("test-data"),
				nil,
			},
			"",
			true,
		},
		{
			"err-nil-msg",
			fields{
				func(ctx aws.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(t, bytes.NewReader([]byte("test-data")), input.Body) {
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
				[]byte("contents-hash"),
				nil,
				nil,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UploadProvider{
				Uploader: tt.fields.uploader,
				Bucket:   tt.fields.bucket,
			}
			key := encoding.EncodeHex(tt.args.contentsHash)
			location, err := u.Upload(context.Background(), tt.args.headers, key, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadProvider.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if location != tt.wantAddress {
				t.Errorf("UploadProvider.Upload() address = %v, wantAddress %v", location, tt.wantAddress)
			}
		})
	}
}
