package address

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_encodeZeroX(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		in []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"success",
			args{
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
			},
			"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeZeroX(tt.args.in); !assert.Equal(tt.want, got) {
				t.Errorf("encodeZeroX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeZeroX(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success",
			args{
				"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
			},
			testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
			false,
		},
		{
			"err-missing-prefix",
			args{
				"5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
			},
			nil,
			true,
		},
		{
			"err-empty",
			args{
				"",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeZeroX(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeZeroX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeZeroX() = %v, want %v", got, tt.want)
			}
		})
	}
}
