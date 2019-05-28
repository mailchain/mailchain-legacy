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

//go:generate mockgen -source=sent.go -package=mocks -destination=$PACKAGE_PATH/internal/testutil/mocks/sent.go

package stores_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/stores"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/mailchain/mailchain/internal/testutil/mocks"
	"github.com/pkg/errors"
)

func TestPutMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		sent      stores.Sent
		messageID mail.ID
		msg       io.Reader
	}
	tests := []struct {
		name         string
		args         args
		wantLocation string
		wantErr      bool
	}{
		{
			"err-nil-sent",
			args{
				nil,
				nil,
				nil,
			},
			"",
			true,
		},
		{
			"err-nil-reader",
			args{
				func() stores.Sent {
					sent := mocks.NewMockSent(mockCtrl)
					return sent
				}(),
				nil,
				nil,
			},
			"",
			true,
		},
		{
			"err-put-error",
			args{
				func() stores.Sent {
					sent := mocks.NewMockSent(mockCtrl)
					sent.EXPECT().PutMessage(
						"002c47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471-22049eeebdc1",
						gomock.Any(),
						nil,
					).Return("", errors.Errorf("put error")).Times(1)
					return sent
				}(),
				testutil.MustHexDecodeString("002c47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"),
				func() io.Reader {
					contents, _ := ioutil.ReadFile("./testdata/simple.golden.eml-22049eeebdc1")
					return bytes.NewReader(contents)
				}(),
			},
			"",
			true,
		},
		{
			"success",
			args{
				func() stores.Sent {
					sent := mocks.NewMockSent(mockCtrl)
					sent.EXPECT().PutMessage(
						"002c47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471-22049eeebdc1",
						gomock.Any(),
						nil,
					).Return("https://location", nil).Times(1)
					return sent
				}(),
				testutil.MustHexDecodeString("002c47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"),
				func() io.Reader {
					contents, _ := ioutil.ReadFile("./testdata/simple.golden.eml-22049eeebdc1")
					return bytes.NewReader(contents)
				}(),
			},
			"https://location",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLocation, err := stores.PutMessage(tt.args.sent, tt.args.messageID, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLocation != tt.wantLocation {
				t.Errorf("PutMessage() = %v, want %v", gotLocation, tt.wantLocation)
			}
		})
	}
}
