package storage

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/s3store"
	"github.com/pkg/errors"
)

type S3Store struct {
	headObjectFunc func(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
	sent           stores.Sent
	bucket         string
}

func (s S3Store) Exists(messageID mail.ID, contents []byte, hash string) error {
	_, err := s.headObjectFunc(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.sent.Key(messageID, contents)),
	})
	if err == nil {
		return errors.Errorf("message already exists")
	}
	aerr, ok := err.(awserr.Error)
	if !ok {
		return err
	}
	if aerr.Code() != "NotFound" {
		return aerr
	}
	return nil
}

func (s S3Store) Put(messageID mail.ID, contents []byte, hash string) (string, error) {
	loc := s.sent.Key(messageID, contents)
	putLoc, err := s.sent.PutMessage(messageID, contents, nil)
	if err != nil {
		return "", errors.WithMessage(err, "could not PUT message")
	}
	if !strings.HasSuffix(putLoc, loc) || strings.TrimSpace(loc) == "" {
		return "", errors.Errorf("message location could not be safely determined")
	}

	// if !strings.HasSuffix(loc, k) {
	// 	return "", errors.Errorf("object not created at correct location")
	// }
	return loc, nil
}

func createS3Client(region, id, secret string) (*s3.S3, error) {
	if id != "" && secret != "" {
		creds := credentials.NewStaticCredentials(id, secret, "")
		ses, err := session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		})
		return s3.New(ses), err
	}

	return s3.New(session.New()), nil
}

// NewSentStore creates a new S3 store.
func NewSentStore(region, bucket, id, secret string) (*S3Store, error) {
	if region == "" {
		return nil, errors.Errorf("`region` must be specified")
	}
	if bucket == "" {
		return nil, errors.Errorf("`bucket` must be specified")
	}
	s3Client, err := createS3Client(region, id, secret)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create session")
	}

	sent, err := s3store.NewSent(region, bucket, id, secret)
	// err handled as part of return

	return &S3Store{
		headObjectFunc: s3Client.HeadObject,
		sent:           sent,
		bucket:         bucket,
	}, errors.WithMessage(err, "could not sent store")
}
