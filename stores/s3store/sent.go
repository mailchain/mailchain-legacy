// Copyright 2022 Mailchain Ltd.
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
	"context"

	"github.com/mailchain/mailchain/encoding"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/internal/mail"
)

// NewSent creates a new S3 store.
func NewSent(region, bucket, id, secret string) (*Sent, error) {
	s3Store, err := NewUploader(region, bucket, id, secret)
	if err != nil {
		return nil, err
	}

	return &Sent{uploader: s3Store}, nil
}

// Sent handles storing messages in S3
type Sent struct {
	uploader Uploader
}

// Key of resource stored.
func (h Sent) Key(messageID mail.ID, contentsHash, msg []byte) string {
	return encoding.EncodeHex(contentsHash)
}

// PutMessage stores the message in S3.
func (h Sent) PutMessage(messageID mail.ID, contentsHash, msg []byte, headers map[string]string) (
	address, resource string, mli uint64, err error) {
	metadata := map[string]*string{
		"Version": aws.String(mailchain.Version),
	}

	for k, v := range headers {
		metadata[k] = aws.String(v)
	}
	resource = h.Key(messageID, contentsHash, msg)

	location, err := h.uploader.Upload(context.Background(), metadata, resource, msg)
	if err != nil {
		return "", "", 0, err
	}

	return location, resource, 0, nil
}
