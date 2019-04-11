// Copyright (c) 2019 Finobo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package rfc2822

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/mailchain/mailchain/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update .golden files")

func TestEncodeNewMessage(t *testing.T) {
	assert := assert.New(t)
	type args struct{ message *mail.Message }
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple",
			args{&mail.Message{
				Headers: &mail.Headers{
					Date:    time.Date(2019, 3, 12, 20, 23, 13, 45, time.UTC),
					From:    mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: ""},
					To:      mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: ""},
					Subject: "Hello world",
				},
				ID:   mail.ID(testutil.MustHexDecodeMultiHashID("002c47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")),
				Body: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Phasellus eget euismod nibh. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer bibendum venenatis sem sed auctor. Ut aliquam eu diam nec fermentum. Sed turpis nulla, viverra ac efficitur ac, fermentum vel sapien. Curabitur vehicula risus id odio congue tempor. Mauris tincidunt feugiat risus, eget auctor magna blandit sit amet. Curabitur consectetur, dolor eu imperdiet varius, dui neque mattis neque, vel fringilla magna tortor ut risus. Cras cursus sem et nisl interdum molestie. Aliquam auctor sodales blandit."),
			}},
			false,
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
				ioutil.WriteFile(golden, got, 0644)
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
				ID:   mail.ID(testutil.MustHexDecodeMultiHashID("002c47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")),
				Headers: &mail.Headers{
					Date:    time.Date(2019, 3, 12, 20, 23, 13, 0, time.UTC),
					From:    mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
					To:      mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
					ReplyTo: nil,
					Subject: "Hello world",
				},
				// TODO: publicKey?
			},
			false,
		},
		// TODO: Add test cases.
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
