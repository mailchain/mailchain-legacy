package substraterpc

import (
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	server := httptest.NewServer(nil)
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			args{
				server.URL,
			},
			false,
			false,
		},
		{
			"failed",
			args{
				"host:23425",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("New() got = %v, wantNil %v", err, tt.wantErr)
				return
			}
		})
	}
}
