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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

// NewSenderStore creates a new S3 store.
func NewSenderStore(region, bucket, id, secret string) (*SenderStore, error) {
	if region == "" {
		return nil, errors.Errorf("`region` must be specified")
	}
	if bucket == "" {
		return nil, errors.Errorf("`bucket` must be specified")
	}
	var creds *credentials.Credentials
	if id != "" && secret != "" {
		creds = credentials.NewStaticCredentials(id, secret, "")
	}

	ses, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "could not create session")
	}

	// S3 service client the Upload manager will use.
	s3Svc := s3.New(ses)

	return &SenderStore{
		uploader: s3manager.NewUploaderWithClient(s3Svc), // Create an uploader with S3 client and default options
		s3:       s3Svc,
		bucket:   bucket,
	}, err
}

// SenderStore handles storing messages in S3
type SenderStore struct {
	uploader *s3manager.Uploader
	s3       *s3.S3
	bucket   string
}

func (h SenderStore) PutMessage(path string, msg []byte) (string, error) {
	upParams := &s3manager.UploadInput{
		Bucket: &h.bucket,
		Key:    &path,
		Body:   bytes.NewReader(msg),
	}
	// Perform an upload.
	result, err := h.uploader.Upload(upParams)
	if err != nil {
		return "", errors.WithMessage(err, "could not put message")
	}
	return result.Location, nil
}
