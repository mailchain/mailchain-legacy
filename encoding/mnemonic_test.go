package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeMnemonicAlgorand(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-25-words",
			args{
				"success kind profit hamster middle drama crouch cat glass sea warm song coyote vacant sport sentence soul decorate shuffle blame unveil snack swim abandon concert",
			},
			[]byte{0xc2, 0xa6, 0x9e, 0x57, 0x8b, 0x26, 0xc6, 0x8, 0x85, 0xc6, 0x23, 0x17, 0x7b, 0x70, 0xee, 0xf3, 0xdc, 0x98, 0xc1, 0x57, 0xfa, 0xc3, 0x7d, 0x46, 0xce, 0x8e, 0x73, 0x31, 0x77, 0x34, 0x7f, 0x1b},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeMnemonicAlgorand(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeMnemonicAlgorand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("DecodeMnemonicAlgorand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeMnemonicAlgorand(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"success-32-bytes",
			args{
				[]byte{0xc2, 0xa6, 0x9e, 0x57, 0x8b, 0x26, 0xc6, 0x8, 0x85, 0xc6, 0x23, 0x17, 0x7b, 0x70, 0xee, 0xf3, 0xdc, 0x98, 0xc1, 0x57, 0xfa, 0xc3, 0x7d, 0x46, 0xce, 0x8e, 0x73, 0x31, 0x77, 0x34, 0x7f, 0x1b},
			},
			"success kind profit hamster middle drama crouch cat glass sea warm song coyote vacant sport sentence soul decorate shuffle blame unveil snack swim abandon concert",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeMnemonicAlgorand(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeMnemonicAlgorand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeMnemonicAlgorand() = %v, want %v", got, tt.want)
			}
		})
	}
}
