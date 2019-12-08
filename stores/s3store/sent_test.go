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

package s3store

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewSent(t *testing.T) {
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
			got, err := NewSent(tt.args.region, tt.args.bucket, tt.args.id, tt.args.secret)
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

func TestSent_PutMessage(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		uploader func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
		bucket   string
	}
	type args struct {
		messageID    mail.ID
		contentsHash []byte
		msg          []byte
		headers      map[string]string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAddress  string
		wantResource string
		wantMLI      uint64
		wantErr      bool
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
					if !assert.Equal(aws.String("636f6e74656e74732d68617368"), input.Key) {
						t.Errorf("Bucket incorrect")
					}

					return &s3manager.UploadOutput{Location: "https://bucket-id/636f6e74656e74732d68617368"}, nil
				},
				"bucket-id",
			},
			args{
				func() mail.ID {
					return []byte("location")
				}(),
				[]byte("contents-hash"),
				[]byte("test-data"),
				nil,
			},
			"https://bucket-id/636f6e74656e74732d68617368",
			"636f6e74656e74732d68617368",
			0,
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
					if !assert.Equal(aws.String("636f6e74656e74732d68617368"), input.Key) {
						t.Errorf("Key incorrect")
					}

					return &s3manager.UploadOutput{Location: "https://bucket-id/636f6e74656e74732d68617368"}, nil
				},
				"bucket-id",
			},
			args{
				func() mail.ID {
					return []byte("location")
				}(),
				[]byte("contents-hash"),
				[]byte("test-data"),
				map[string]string{
					"key-1": "value-1",
				},
			},
			"https://bucket-id/636f6e74656e74732d68617368",
			"636f6e74656e74732d68617368",
			0,
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
					if !assert.Equal(aws.String("636f6e74656e74732d68617368"), input.Key) {
						t.Errorf("Key incorrect")
					}

					return nil, errors.Errorf("failed to upload")
				},
				"bucket-id",
			},
			args{
				func() mail.ID {
					return []byte("location")
				}(),
				[]byte("contents-hash"),
				[]byte("test-data"),
				nil,
			},
			"",
			"",
			0,
			true,
		},
		{
			"err-nil-msg",
			fields{
				func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
					if !assert.Equal(bytes.NewReader([]byte("test-data")), input.Body) {
						t.Errorf("body incorrect")
					}
					if !assert.Equal(aws.String("bucket-id"), input.Bucket) {
						t.Errorf("Bucket incorrect")
					}
					if !assert.Equal(aws.String("636f6e74656e74732d68617368"), input.Key) {
						t.Errorf("Key incorrect")
					}

					return nil, errors.Errorf("failed to upload")
				},
				"bucket-id",
			},
			args{
				func() mail.ID {
					return []byte("location")
				}(),
				[]byte("contents-hash"),
				nil,
				nil,
			},
			"",
			"",
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Sent{
				uploader: tt.fields.uploader,
				bucket:   tt.fields.bucket,
			}
			gotAddress, gotResource, gotMLI, err := h.PutMessage(tt.args.messageID, tt.args.contentsHash, tt.args.msg, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sent.PutMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("Sent.PutMessage() address = %v, wantAddress %v", gotAddress, tt.wantAddress)
			}
			if gotResource != tt.wantResource {
				t.Errorf("Sent.PutMessage() resource = %v, wantResource %v", gotResource, tt.wantResource)
			}
			if gotMLI != tt.wantMLI {
				t.Errorf("Sent.PutMessage() = %v, wantMLI %v", gotMLI, tt.wantMLI)
			}
		})
	}
}

func TestSent_Key(t *testing.T) {
	type fields struct {
		uploader func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
		bucket   string
	}
	type args struct {
		messageID    mail.ID
		contentsHash []byte
		msg          []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"success",
			fields{
				nil,
				"",
			},
			args{
				[]byte("messageID"),
				[]byte("contents-hash"),
				[]byte("body"),
			},
			encoding.EncodeHex([]byte("contents-hash")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Sent{
				uploader: tt.fields.uploader,
				bucket:   tt.fields.bucket,
			}
			if got := h.Key(tt.args.messageID, tt.args.contentsHash, tt.args.msg); got != tt.want {
				t.Errorf("Sent.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}
