// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewSentStore(t *testing.T) {
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
				"bucket",
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
				"bucket",
				"id",
				"secret",
			},
			true,
			true,
		},
		{
			"err-bucket",
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
			got, err := NewSentStore(tt.args.region, tt.args.bucket, tt.args.id, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSentStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil != tt.wantNil {
				t.Errorf("NewSentStore() got == nil = %v, wantNil %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestSentStore_PutMessage(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		uploader func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
		bucket   string
	}
	type args struct {
		path    string
		msg     []byte
		headers map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"success-no-headers",
			fields{
				func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(bytes.NewReader([]byte("test-data")), input.Body) {
						t.Errorf("body incorrect")
					}
					if !assert.Equal(aws.String("bucket-id"), input.Bucket) {
						t.Errorf("Bucket incorrect")
					}
					if !assert.Equal(aws.String("location-hash"), input.Key) {
						t.Errorf("Key incorrect")
					}

					return &s3manager.UploadOutput{Location: "https://bucket-id/location-hash"}, nil
				},
				"bucket-id",
			},
			args{
				"location-hash",
				[]byte("test-data"),
				nil,
			},
			"https://bucket-id/location-hash",
			false,
		},
		{
			"success-has-headers",
			fields{
				func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(bytes.NewReader([]byte("test-data")), input.Body) {
						t.Errorf("body incorrect")
					}
					if !assert.Equal(aws.String("bucket-id"), input.Bucket) {
						t.Errorf("Bucket incorrect")
					}
					if !assert.Equal(aws.String("location-hash"), input.Key) {
						t.Errorf("Key incorrect")
					}

					return &s3manager.UploadOutput{Location: "https://bucket-id/location-hash"}, nil
				},
				"bucket-id",
			},
			args{
				"location-hash",
				[]byte("test-data"),
				map[string]string{
					"key-1": "value-1",
				},
			},
			"https://bucket-id/location-hash",
			false,
		},
		{
			"err-uploader",
			fields{
				func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(bytes.NewReader([]byte("test-data")), input.Body) {
						t.Errorf("body incorrect")
					}
					if !assert.Equal(aws.String("bucket-id"), input.Bucket) {
						t.Errorf("Bucket incorrect")
					}
					if !assert.Equal(aws.String("location-hash"), input.Key) {
						t.Errorf("Key incorrect")
					}

					return nil, errors.Errorf("failed to upload")
				},
				"bucket-id",
			},
			args{
				"location-hash",
				[]byte("test-data"),
				nil,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := SentStore{
				uploader: tt.fields.uploader,
				bucket:   tt.fields.bucket,
			}
			got, err := h.PutMessage(tt.args.path, tt.args.msg, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentStore.PutMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SentStore.PutMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
