package address

import (
	"testing"

	"github.com/mailchain/mailchain/internal/testutil"
)

func TestEncodeByProtocol(t *testing.T) {
	type args struct {
		in       []byte
		protocol string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"ethereum",
			args{
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				"ethereum",
			},
			"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
			false,
		},
		{
			"err",
			args{
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				"invalid",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeByProtocol(tt.args.in, tt.args.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeByProtocol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeByProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}
