package sr25519

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/internal/testutil"
	//"github.com/stretchr/testify/assert"
)

func TestPrivateKeyFromBytes(t *testing.T) {
	type args struct {
		pk []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *PrivateKey
		wantErr bool
	}{
		{
			"success-sofia-seed",
			args{
				sofiaSeed,
			},
			&sofiaPrivateKey,
			false,
		},
		{
			"success-sofia-bytes",
			args{
				sofiaPrivateKeyBytes,
			},
			&sofiaPrivateKey,
			false,
		},
		{
			"success-charlotte-seed",
			args{
				charlotteSeed,
			},
			&charlottePrivateKey,
			false,
		},
		{
			"success-charlotte-bytes",
			args{
				charlottePrivateKeyBytes,
			},
			&charlottePrivateKey,
			false,
		},
		{
			"err-len",
			args{
				testutil.MustHexDecodeString("39d4c9"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyFromBytes(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
