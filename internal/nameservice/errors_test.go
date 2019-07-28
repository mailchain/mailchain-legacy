package nameservice

import (
	"testing"

	"github.com/pkg/errors"
)

func TestIsNoResolverError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"ErrUnableToResolve",
			args{
				errors.Errorf("unable to resolve and other"),
			},
			true,
		},
		{
			"other",
			args{
				errors.Errorf("other"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNoResolverError(tt.args.err); got != tt.want {
				t.Errorf("IsNoResolverError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"ErrNotFound",
			args{
				errors.Errorf("not found and other"),
			},
			true,
		},
		{
			"other",
			args{
				errors.Errorf("other"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFoundError(tt.args.err); got != tt.want {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvalidNameError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"ErrInvalidName",
			args{
				errors.Errorf("invalid name"),
			},
			true,
		},
		{
			"other",
			args{
				errors.Errorf("other"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInvalidNameError(tt.args.err); got != tt.want {
				t.Errorf("IsInvalidNameError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvalidAddressError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"ErrInvalidAddress",
			args{
				errors.Errorf("invalid address"),
			},
			true,
		},
		{
			"other",
			args{
				errors.Errorf("other"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInvalidAddressError(tt.args.err); got != tt.want {
				t.Errorf("IsInvalidAddressError() = %v, want %v", got, tt.want)
			}
		})
	}
}
