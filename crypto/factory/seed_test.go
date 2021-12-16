package factory

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func Test_generateSeed(t *testing.T) {
	type args struct {
		r      io.Reader
		length uint8
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
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
				16,
			},
			[]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50},
			false,
		},
		{
			"err",
			args{
				iotest.ErrReader(errors.New("error")),
				255,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSeed(tt.args.r, tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("generateSeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
