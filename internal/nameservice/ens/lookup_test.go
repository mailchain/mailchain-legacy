package ens

import (
	"testing"
)

func TestNewLookupService(t *testing.T) {
	type args struct {
		clientURL string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			"success",
			args{
				"https://client.url",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLookupService(tt.args.clientURL)
			if (got == nil) != tt.wantNil {
				t.Errorf("NewLookupService() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}
