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
	nm "net/mail"
	"testing"
	"time"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/stretchr/testify/assert"
)

func Test_parseSubject(t *testing.T) {
	type args struct {
		h nm.Header
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"success",
			args{
				nm.Header{
					"Subject": []string{"test subject"},
				},
			},
			"test subject",
			false,
		},
		{
			"err-empty",
			args{
				nm.Header{
					"Subject": []string{},
				},
			},
			"",
			true,
		},
		{
			"err-missing",
			args{
				nm.Header{},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSubject(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSubject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDate(t *testing.T) {
	type args struct {
		h nm.Header
	}
	tests := []struct {
		name    string
		args    args
		want    *time.Time
		wantErr bool
	}{
		{
			"success",
			args{
				nm.Header{
					"Date": []string{"Tue, 12 Mar 2019 20:23:13 UTC"},
				},
			},
			func() *time.Time {
				t := time.Date(2019, 03, 12, 20, 23, 13, 0, time.UTC)
				return &t
			}(),
			false,
		},
		{
			"err-invalid",
			args{
				nm.Header{
					"Date": []string{"Tue, 32 Mar 2019 20:23:13 UTC"},
				},
			},
			nil,
			true,
		},
		{
			"err-empty",
			args{
				nm.Header{
					"Date": []string{},
				},
			},
			nil,
			true,
		},
		{
			"err-missing",
			args{
				nm.Header{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFrom(t *testing.T) {
	type args struct {
		h nm.Header
	}
	tests := []struct {
		name    string
		args    args
		want    *mail.Address
		wantErr bool
	}{
		{
			"success",
			args{
				nm.Header{
					"From": []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
				},
			},
			&mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
			false,
		},
		{
			"err-empty",
			args{
				nm.Header{
					"From": []string{},
				},
			},
			nil,
			true,
		},
		{
			"err-missing",
			args{
				nm.Header{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFrom(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTo(t *testing.T) {
	type args struct {
		h nm.Header
	}
	tests := []struct {
		name    string
		args    args
		want    *mail.Address
		wantErr bool
	}{
		{
			"success",
			args{
				nm.Header{
					"To": []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
				},
			},
			&mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
			false,
		},
		{
			"err-empty",
			args{
				nm.Header{
					"To": []string{},
				},
			},
			nil,
			true,
		},
		{
			"err-missing",
			args{
				nm.Header{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTo(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseContentType(t *testing.T) {

	type args struct {
		h nm.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"success-plain-text",
			args{
				nm.Header{
					"Content-Type": []string{"text/plain; charset=\"UTF-8\""},
				},
			},
			"text/plain; charset=\"UTF-8\"",
		},
		{
			"success-empty-header-to-default",
			args{
				nm.Header{},
			},
			"text/plain; charset=\"UTF-8\"",
		},
		{
			"success-empty-header-value-to-default",
			args{
				nm.Header{
					"Content-Type": []string{""},
				},
			},
			"text/plain; charset=\"UTF-8\"",
		},
		{
			"success-html-text",
			args{
				nm.Header{
					"Content-Type": []string{"text/html; charset=\"UTF-8\""},
				},
			},
			"text/html; charset=\"UTF-8\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseContentType(tt.args.h)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePublicKey(t *testing.T) {
	pubKeyHeader := func(pk crypto.PublicKey) string {
		keyKind, err := multikey.KindFromPublicKey(pk)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		return encoding.EncodeHexZeroX(pk.Bytes()) + "; type=\"" + keyKind + "\"; encoding=" + encoding.KindHex0XPrefix
	}
	type args struct {
		h nm.Header
	}
	tests := []struct {
		name    string
		args    args
		want    crypto.PublicKey
		wantErr bool
	}{
		{
			"secp256k1",
			args{
				nm.Header{
					"Public-Key": []string{pubKeyHeader(secp256k1test.AlicePublicKey)},
				},
			},
			secp256k1test.AlicePublicKey,
			false,
		},
		{
			"ed25519",
			args{
				nm.Header{
					"Public-Key": []string{pubKeyHeader(ed25519test.AlicePublicKey)},
				},
			},
			ed25519test.AlicePublicKey,
			false,
		},
		{
			"missing",
			args{
				nm.Header{},
			},
			nil,
			false,
		},
		{
			"err-missing-delimiter",
			args{
				nm.Header{
					"Public-Key": []string{encoding.EncodeHexZeroX(ed25519test.AlicePublicKey.Bytes()) + " type=ed25519 encoding=" + encoding.KindHex0XPrefix},
				},
			},
			nil,
			true,
		},
		{
			"err-invalid-key-bytes",
			args{
				nm.Header{
					"Public-Key": []string{"invalid-key-bytes; type=ed25519 ;encoding=" + encoding.KindHex0XPrefix},
				},
			},
			nil,
			true,
		},
		{
			"err-missing-type",
			args{
				nm.Header{
					"Public-Key": []string{encoding.EncodeHexZeroX(ed25519test.AlicePublicKey.Bytes()) + "; encoding=" + encoding.KindHex0XPrefix},
				},
			},
			nil,
			true,
		},
		{
			"err-empty",
			args{
				nm.Header{
					"Public-Key": []string{},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePublicKey(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePublicKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !assert.EqualValues(t, tt.want, got) {
				t.Errorf("PublicKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHeaders(t *testing.T) {
	pubKeyHeader := func(pk crypto.PublicKey) string {
		keyKind, err := multikey.KindFromPublicKey(pk)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		return encoding.EncodeHexZeroX(pk.Bytes()) + "; type=" + keyKind + "; encoding=" + encoding.KindHex0XPrefix
	}

	type args struct {
		h nm.Header
	}
	tests := []struct {
		name    string
		args    args
		want    *mail.Headers
		wantErr bool
	}{
		{
			"success-plain-text",
			args{
				nm.Header{
					"Subject":      []string{"test subject"},
					"To":           []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
					"From":         []string{"<4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>"},
					"Date":         []string{"Tue, 12 Mar 2019 20:23:13 UTC"},
					"Content-Type": []string{"text/plain; charset=\"UTF-8\""},
					"Public-Key":   []string{pubKeyHeader(ed25519test.AlicePublicKey)},
				},
			},
			&mail.Headers{
				From:        mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
				To:          mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				Date:        time.Date(2019, 03, 12, 20, 23, 13, 0, time.UTC),
				Subject:     "test subject",
				ReplyTo:     nil,
				ContentType: "text/plain; charset=\"UTF-8\"",
				PublicKey:   ed25519test.AlicePublicKey,
			},
			false,
		},
		{
			"success-plain-html",
			args{
				nm.Header{
					"Subject":      []string{"test subject"},
					"To":           []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
					"From":         []string{"<4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>"},
					"Date":         []string{"Tue, 12 Mar 2019 20:23:13 UTC"},
					"Content-Type": []string{"text/html; charset=\"UTF-8\""},
					"Public-Key":   []string{pubKeyHeader(secp256k1test.AlicePublicKey)},
				},
			},
			&mail.Headers{
				From:        mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
				To:          mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				Date:        time.Date(2019, 03, 12, 20, 23, 13, 0, time.UTC),
				Subject:     "test subject",
				ReplyTo:     nil,
				ContentType: "text/html; charset=\"UTF-8\"",
				PublicKey:   secp256k1test.AlicePublicKey,
			},
			false,
		},
		{
			"success-defaultContentType",
			args{
				nm.Header{
					"Subject":    []string{"test subject"},
					"To":         []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
					"From":       []string{"<4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>"},
					"Date":       []string{"Tue, 12 Mar 2019 20:23:13 UTC"},
					"Public-Key": []string{pubKeyHeader(secp256k1test.AlicePublicKey)},
				},
			},
			&mail.Headers{
				From:        mail.Address{DisplayName: "", FullAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum", ChainAddress: "4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"},
				To:          mail.Address{DisplayName: "", FullAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum", ChainAddress: "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"},
				Date:        time.Date(2019, 03, 12, 20, 23, 13, 0, time.UTC),
				Subject:     "test subject",
				ReplyTo:     nil,
				ContentType: "text/plain; charset=\"UTF-8\"",
				PublicKey:   secp256k1test.AlicePublicKey,
			},
			false,
		},
		{
			"err-from",
			args{
				nm.Header{
					"Subject":    []string{"test subject"},
					"To":         []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
					"Date":       []string{"Tue, 12 Mar 2019 20:23:13 UTC"},
					"Public-Key": []string{pubKeyHeader(secp256k1test.AlicePublicKey)},
				},
			},
			nil,
			true,
		},
		{
			"err-to",
			args{
				nm.Header{
					"Subject":    []string{"test subject"},
					"From":       []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
					"Date":       []string{"Tue, 12 Mar 2019 20:23:13 UTC"},
					"Public-Key": []string{pubKeyHeader(secp256k1test.AlicePublicKey)},
				},
			},
			nil,
			true,
		},
		{
			"err-subject",
			args{
				nm.Header{
					"To":              []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
					"From":            []string{"<4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>"},
					"Date":            []string{"Tue, 12 Mar 2019 20:23:13 UTC"},
					"Public-Key":      []string{encoding.EncodeHexZeroX(ed25519test.AlicePublicKey.Bytes())},
					"Public-Key-Type": []string{"ed25519"},
				},
			},
			nil,
			true,
		},
		{
			"err-date",
			args{
				nm.Header{
					"Subject":    []string{"test subject"},
					"To":         []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
					"From":       []string{"<4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>"},
					"Public-Key": []string{pubKeyHeader(secp256k1test.AlicePublicKey)},
				},
			},
			nil,
			true,
		},
		{
			"err-public-key",
			args{
				nm.Header{
					"Subject":    []string{"test subject"},
					"To":         []string{"<5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>"},
					"From":       []string{"<4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>"},
					"Date":       []string{"Tue, 12 Mar 2019 20:23:13 UTC"},
					"Public-Key": []string{"invalid"},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeaders(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
