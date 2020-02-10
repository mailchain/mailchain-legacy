package s3store

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mailchain/mailchain/encoding"
	"github.com/pkg/errors"
)

type S3Store struct {
	Uploader func(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	Bucket   string
}

// Key of resource stored.
func (s *S3Store) EncodeKey(hash []byte) string {
	return encoding.EncodeHex(hash)
}

// NewSent creates a new S3 store.
func NewS3Store(region, bucket, id, secret string) (*S3Store, error) {
	if region == "" {
		return nil, errors.Errorf("`region` must be specified")
	}
	if bucket == "" {
		return nil, errors.Errorf("`Bucket` must be specified")
	}

	var creds *credentials.Credentials

	if id != "" && secret != "" {
		creds = credentials.NewStaticCredentials(id, secret, "")
	}
	ses := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	}))

	// S3 service client the Upload manager will use.
	return &S3Store{
		Uploader: s3manager.NewUploaderWithClient(s3.New(ses)).UploadWithContext, // Create an Uploader with S3 client and default options
		Bucket:   bucket,
	}, nil
}

func (s *S3Store) Upload(ctx context.Context, metadata map[string]*string, key string, msg []byte) (address, resource string, mli uint64, err error) {
	if msg == nil {
		return "", "", 0, errors.Errorf("'msg' must not be nil")
	}

	params := &s3manager.UploadInput{
		Bucket:   &s.Bucket,
		Key:      aws.String(key),
		Body:     bytes.NewReader(msg),
		Metadata: metadata,
	}
	// Perform an upload.
	result, err := s.Uploader(ctx, params)
	if err != nil {
		return "", "", 0, errors.WithMessage(err, "could not put message")
	}
	return result.Location, key, 0, nil
}
