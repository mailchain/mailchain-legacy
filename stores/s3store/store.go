package s3store

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

type Uploader struct {
	Uploader func(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	Bucket   string
}

// NewUploader creates a new Uploader store for S3.
func NewUploader(region, bucket, id, secret string) (*Uploader, error) {
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
	return &Uploader{
		Uploader: s3manager.NewUploaderWithClient(s3.New(ses)).UploadWithContext, // Create an Uploader with S3 client and default options
		Bucket:   bucket,
	}, nil
}

func (s *Uploader) Upload(ctx context.Context, metadata map[string]*string, key string, msg []byte) (string, error) {
	if msg == nil {
		return "", errors.Errorf("'msg' must not be nil")
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
		return "", errors.WithMessage(err, "could not put message")
	}

	return result.Location, nil
}
