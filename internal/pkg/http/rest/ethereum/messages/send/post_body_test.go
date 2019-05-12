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

package send

import (
	"encoding/hex"
	"testing"

	"github.com/mailchain/mailchain/internal/pkg/testutil"
)

func Test_checkForEmpties(t *testing.T) {
	type args struct {
		msg PostMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{
				PostMessage{
					Headers:   &PostHeaders{},
					Subject:   "subject-value",
					Body:      "body-value",
					PublicKey: "public-key-value",
				},
			},
			false,
		},
		{
			"empty-headers",
			args{
				PostMessage{
					Subject:   "subject-value",
					Body:      "body-value",
					PublicKey: "public-key-value",
				},
			},
			true,
		},
		{
			"empty-subject",
			args{
				PostMessage{
					Headers:   &PostHeaders{},
					Body:      "body-value",
					PublicKey: "public-key-value",
				},
			},
			true,
		},
		{
			"empty-body",
			args{
				PostMessage{
					Headers:   &PostHeaders{},
					Subject:   "subject-value",
					PublicKey: "public-key-value",
				},
			},
			true,
		},
		{
			"empty-public-key",
			args{
				PostMessage{
					Headers: &PostHeaders{},
					Subject: "subject-value",
					Body:    "body-value",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkForEmpties(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("checkForEmpties() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isValid(t *testing.T) {
	type args struct {
		p       *PostRequestBody
		network string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"err-nil",
			args{},
			true,
		},
		{
			"err-empties",
			args{
				&PostRequestBody{},
				"ethereum",
			},
			true,
		},
		{
			"err-parse-to",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers:   &PostHeaders{},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "public-key-value",
					},
				},
				"ethereum",
			},
			true,
		},
		{
			"err-parse-from",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To: hex.EncodeToString(testutil.CharlottePublicKey.Address()),
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "public-key-value",
					},
				},
				"ethereum",
			},
			true,
		},
		{
			"err-reply-to",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To:      hex.EncodeToString(testutil.CharlottePublicKey.Address()),
							From:    hex.EncodeToString(testutil.SofiaPublicKey.Address()),
							ReplyTo: "<invalid",
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "public-key-value",
					},
				},
				"ethereum",
			},
			true,
		},
		{
			"err-public-key",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To:   hex.EncodeToString(testutil.CharlottePublicKey.Address()),
							From: hex.EncodeToString(testutil.SofiaPublicKey.Address()),
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "public-key-value",
					},
				},
				"ethereum",
			},
			true,
		},
		{
			"err-address-from-public-key",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To:   hex.EncodeToString(testutil.CharlottePublicKey.Address()),
							From: hex.EncodeToString(testutil.SofiaPublicKey.Address()),
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "0x" + hex.EncodeToString(testutil.SofiaPublicKey.Bytes()),
					},
				},
				"ethereum",
			},
			true,
		},
		{
			"success",
			args{
				&PostRequestBody{
					Message: PostMessage{
						Headers: &PostHeaders{
							To:   hex.EncodeToString(testutil.CharlottePublicKey.Address()),
							From: hex.EncodeToString(testutil.SofiaPublicKey.Address()),
						},
						Subject:   "subject-value",
						Body:      "body-value",
						PublicKey: "0x" + hex.EncodeToString(testutil.CharlottePublicKey.Bytes()),
					},
				},
				"ethereum",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isValid(tt.args.p, tt.args.network); (err != nil) != tt.wantErr {
				t.Errorf("isValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
