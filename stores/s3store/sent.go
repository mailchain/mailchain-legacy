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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

// NewSent creates a new S3 store.
func NewSent(region, bucket, id, secret string) (*Sent, error) {
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
	ses := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	}))

	// S3 service client the Upload manager will use.
	return &Sent{
		uploader: s3manager.NewUploaderWithClient(s3.New(ses)).Upload, // Create an uploader with S3 client and default options
		bucket:   bucket,
	}, nil
}

// Sent handles storing messages in S3
type Sent struct {
	uploader func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	bucket   string
}

func (h Sent) Key(messageID mail.ID, msg []byte) string {
	hash := crypto.CreateLocationHash(msg)
	return fmt.Sprintf("%s-%s", messageID.HexString(), hash.String())
}

func (h Sent) PutMessage(messageID mail.ID, msg []byte, headers map[string]string) (string, error) {
	if msg == nil {
		return "", errors.Errorf("'msg' must not be nil")
	}
	metadata := map[string]*string{
		"Version": aws.String(mailchain.Version),
	}
	for k, v := range headers {
		metadata[k] = aws.String(v)
	}
	params := &s3manager.UploadInput{
		Bucket:   &h.bucket,
		Key:      aws.String(h.Key(messageID, msg)),
		Body:     bytes.NewReader(msg),
		Metadata: metadata,
	}
	// Perform an upload.
	result, err := h.uploader(params)
	if err != nil {
		return "", errors.WithMessage(err, "could not put message")
	}
	return result.Location, nil
}
