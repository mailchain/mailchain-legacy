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

package rfc2822

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update .golden files")

func TestEncodeNewMessage(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		message *mail.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple",
			args{&mail.Message{
				Headers: &mail.Headers{
					Date:        time.Date(2019, 3, 12, 20, 23, 13, 45, time.UTC),
					From:        mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: ""},
					To:          mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: ""},
					Subject:     "Hello world",
					ContentType: "text/plain; charset=\"UTF-8\"",
				},
				ID:   mail.ID(testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")),
				Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Phasellus eget euismod nibh. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer bibendum venenatis sem sed auctor. Ut aliquam eu diam nec fermentum. Sed turpis nulla, viverra ac efficitur ac, fermentum vel sapien. Curabitur vehicula risus id odio congue tempor. Mauris tincidunt feugiat risus, eget auctor magna blandit sit amet. Curabitur consectetur, dolor eu imperdiet varius, dui neque mattis neque, vel fringilla magna tortor ut risus. Cras cursus sem et nisl interdum molestie. Aliquam auctor sodales blandit."),
			}},
			false,
		},
		{"html",
			args{&mail.Message{
				Headers: &mail.Headers{
					Date:        time.Date(2019, 3, 12, 20, 23, 13, 45, time.UTC),
					From:        mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: ""},
					To:          mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: ""},
					Subject:     "Hello world",
					ContentType: "text/html; charset=\"UTF-8\"",
				},
				ID:   mail.ID(testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")),
				Body: []byte("<html><h1>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</h1><p>Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Phasellus eget euismod nibh. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer bibendum venenatis sem sed auctor. Ut aliquam eu diam nec fermentum. Sed turpis nulla, viverra ac efficitur ac, fermentum vel sapien. Curabitur vehicula risus id odio congue tempor. Mauris tincidunt feugiat risus, eget auctor magna blandit sit amet. Curabitur consectetur, dolor eu imperdiet varius, dui neque mattis neque, vel fringilla magna tortor ut risus. Cras cursus sem et nisl interdum molestie. Aliquam auctor sodales blandit.</p></html>"),
			}},
			false,
		},
		{"reply-to",
			args{&mail.Message{
				Headers: &mail.Headers{
					Date:        time.Date(2019, 3, 12, 20, 23, 13, 45, time.UTC),
					From:        mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: ""},
					To:          mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: ""},
					ReplyTo:     &mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@mainnet.ethereum", ChainAddress: ""},
					Subject:     "Hello world",
					ContentType: "text/plain; charset=\"UTF-8\"",
				},
				ID:   mail.ID(testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")),
				Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Phasellus eget euismod nibh. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer bibendum venenatis sem sed auctor. Ut aliquam eu diam nec fermentum. Sed turpis nulla, viverra ac efficitur ac, fermentum vel sapien. Curabitur vehicula risus id odio congue tempor. Mauris tincidunt feugiat risus, eget auctor magna blandit sit amet. Curabitur consectetur, dolor eu imperdiet varius, dui neque mattis neque, vel fringilla magna tortor ut risus. Cras cursus sem et nisl interdum molestie. Aliquam auctor sodales blandit."),
			}},
			false,
		},
		{
			"err-nil-message",
			args{
				message: nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeNewMessage(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeNewMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			golden := filepath.Join("testdata", tt.name+".golden.eml")
			if *update {
				_ = ioutil.WriteFile(golden, got, 0644)
			}
			want, _ := ioutil.ReadFile(golden)

			assert.EqualValues(want, got)
			assert.Equal(len(got), cap(got))
		})
	}
}

func TestDecodeNewMessage(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		want    *mail.Message
		wantErr bool
	}{
		// TODO: display names
		// TODO: reply to
		{
			"simple",
			&mail.Message{
				Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Phasellus eget euismod nibh. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer bibendum venenatis sem sed auctor. Ut aliquam eu diam nec fermentum. Sed turpis nulla, viverra ac efficitur ac, fermentum vel sapien. Curabitur vehicula risus id odio congue tempor. Mauris tincidunt feugiat risus, eget auctor magna blandit sit amet. Curabitur consectetur, dolor eu imperdiet varius, dui neque mattis neque, vel fringilla magna tortor ut risus. Cras cursus sem et nisl interdum molestie. Aliquam auctor sodales blandit.\r\n"),
				ID:   mail.ID(testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")),
				Headers: &mail.Headers{
					Date:        time.Date(2019, 3, 12, 20, 23, 13, 0, time.UTC),
					From:        mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
					To:          mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
					ReplyTo:     nil,
					Subject:     "Hello world",
					ContentType: "text/plain; charset=\"UTF-8\"",
				},
				// TODO: publicKey?
			},
			false,
		},
		{
			"html",
			&mail.Message{
				Body: []byte("<html><h1>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</h1><p>Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Phasellus eget euismod nibh. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer bibendum venenatis sem sed auctor. Ut aliquam eu diam nec fermentum. Sed turpis nulla, viverra ac efficitur ac, fermentum vel sapien. Curabitur vehicula risus id odio congue tempor. Mauris tincidunt feugiat risus, eget auctor magna blandit sit amet. Curabitur consectetur, dolor eu imperdiet varius, dui neque mattis neque, vel fringilla magna tortor ut risus. Cras cursus sem et nisl interdum molestie. Aliquam auctor sodales blandit.</p></html>\r\n"),
				ID:   mail.ID(testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")),
				Headers: &mail.Headers{
					Date:        time.Date(2019, 3, 12, 20, 23, 13, 0, time.UTC),
					From:        mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
					To:          mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
					ReplyTo:     nil,
					Subject:     "Hello world",
					ContentType: "text/html; charset=\"UTF-8\"",
				},
			},
			false,
		},
		{
			"err-header",
			nil,
			true,
		},
		{
			"err-id",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			golden := filepath.Join("testdata", tt.name+".golden.eml")
			source, _ := os.Open(golden)
			got, err := DecodeNewMessage(source)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeNewMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(got, tt.want) {
				t.Errorf("DecodeNewMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
