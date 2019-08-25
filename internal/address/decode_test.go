package address

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/internal/testutil"
)

func TestDecodeByProtocol(t *testing.T) {
	type args struct {
		in       string
		protocol string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"ethereum",
			args{
				"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"ethereum",
			},
			testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
			false,
		},
		{
			"err",
			args{
				"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"invalid",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeByProtocol(tt.args.in, tt.args.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeByProtocol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeByProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}
