package storage

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/storestest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_createS3Client(t *testing.T) {
	type args struct {
		region string
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
			"success-without-credentials",
			args{
				"region",
				"",
				"",
			},
			false,
			false,
		},
		{
			"success-with-credentials",
			args{
				"region",
				"id",
				"secret",
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createS3Client(tt.args.region, tt.args.id, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("createSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantNil, got == nil) {
				t.Errorf("createSession() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

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
			"err-region-empty",
			args{
				"",
				"",
				"",
				"",
			},
			true,
			true,
		},
		{
			"err-bucket-empty",
			args{
				"us-east-1",
				"",
				"",
				"",
			},
			true,
			true,
		},
		{
			"success",
			args{
				"us-east-1",
				"bucket",
				"",
				"",
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSentStore(tt.args.region, tt.args.bucket, tt.args.id, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSentStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantNil, got == nil) {
				t.Errorf("NewSentStore() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func TestS3Store_Put(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		headObjectFunc func(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
		sent           stores.Sent
		bucket         string
	}
	type args struct {
		messageID     mail.ID
		contentsHash  []byte
		integrityHash []byte
		contents      []byte
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
			"success",
			fields{
				nil,
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("hashkey")
					sent.EXPECT().PutMessage(id, []byte("contents-hash"), []byte("body"), nil).Return("https://s3bucket/hashkey", "hashkey", envelope.MLIMailchain, nil)
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			"https://s3bucket/hashkey",
			"hashkey",
			1,
			false,
		},
		{
			"err-put-message",
			fields{
				nil,
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("messageID-hash")
					sent.EXPECT().PutMessage(id, []byte("contents-hash"), []byte("body"), nil).Return("", "", envelope.MLIMailchain, errors.Errorf("put failed"))
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			"",
			"",
			1,
			true,
		},
		{
			"err-empty-key",
			fields{
				nil,
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("")
					sent.EXPECT().PutMessage(id, []byte("contents-hash"), []byte("body"), nil).Return("https://s3bucket/hashkey", "hashkey", envelope.MLIMailchain, nil)
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			"",
			"",
			1,
			true,
		},
		{
			"err-inconsistent-key",
			fields{
				nil,
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("messageIDother-hashother")
					sent.EXPECT().PutMessage(id, []byte("contents-hash"), []byte("body"), nil).Return("https://s3bucket/hashkey", "hashkey", envelope.MLIMailchain, nil)
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			"",
			"",
			1,
			true,
		},
		{
			"err-inconsistent-resource",
			fields{
				nil,
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("hashkey")
					sent.EXPECT().PutMessage(id, []byte("contents-hash"), []byte("body"), nil).Return("https://s3bucket/hashkey", "inconsistent-resource", envelope.MLIMailchain, nil)
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			"",
			"",
			1,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := S3Store{
				headObjectFunc: tt.fields.headObjectFunc,
				sent:           tt.fields.sent,
				bucket:         tt.fields.bucket,
			}
			gotAddress, gotResource, gotMLI, err := s.Put(tt.args.messageID, tt.args.contentsHash, tt.args.integrityHash, tt.args.contents)
			if (err != nil) != tt.wantErr {
				t.Errorf("S3Store.Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("S3Store.Put() Address = %v, wantAddress %v", gotAddress, tt.wantAddress)
			}
			if gotResource != tt.wantResource {
				t.Errorf("S3Store.Put() Resource = %v, wantResource %v", gotResource, tt.wantResource)
			}
			if gotMLI != tt.wantMLI {
				t.Errorf("S3Store.Put() MLI = %v, want %v", gotMLI, tt.wantMLI)
			}
		})
	}
}

func TestS3Store_Exists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		headObjectFunc func(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
		sent           stores.Sent
		bucket         string
	}
	type args struct {
		messageID     mail.ID
		contentsHash  []byte
		integrityHash []byte
		contents      []byte
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
				func(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
					return nil, awserr.New("NotFound", "test error", nil)
				},
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("messageID-hash")
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			false,
		},
		{
			"err-non-aws-err",
			fields{
				func(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
					return nil, errors.Errorf("other error")
				},
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("messageID-hash")
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			true,
		},
		{
			"err-exists",
			fields{
				func(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
					return &s3.HeadObjectOutput{}, nil
				},
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("messageID-hash")
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			true,
		},
		{
			"err-different-aws-err",
			fields{
				func(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
					return nil, awserr.New("Forbidden", "test error", nil)
				},
				func() stores.Sent {
					sent := storestest.NewMockSent(mockCtrl)
					var id mail.ID
					id = encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
					sent.EXPECT().Key(id, []byte("contents-hash"), []byte("body")).Return("messageID-hash")
					return sent
				}(),
				"bucket",
			},
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				[]byte("contents-hash"),
				[]byte("integrity-hash"),
				[]byte("body"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := S3Store{
				headObjectFunc: tt.fields.headObjectFunc,
				sent:           tt.fields.sent,
				bucket:         tt.fields.bucket,
			}
			if err := s.Exists(tt.args.messageID, tt.args.contentsHash, tt.args.integrityHash, tt.args.contents); (err != nil) != tt.wantErr {
				t.Errorf("S3Store.Exists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
