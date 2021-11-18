package secp256k1

import (
	"reflect"
	"testing"
)

func TestNewChildKey(t *testing.T) {
	type args struct {
		key PublicKey
		i   uint32
	}
	tests := []struct {
		name    string
		args    args
		want    PublicKey
		wantErr bool
	}{
		{
			"test",
			alicePublicKeyBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChildKey(tt.args.key, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChildKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChildKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
