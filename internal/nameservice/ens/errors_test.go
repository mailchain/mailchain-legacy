package ens

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_wrapError(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		err error
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantErrMessage string
	}{
		{
			"nil",
			args{
				nil,
			},
			false,
			"",
		},
		{
			"no-resolver",
			args{
				errors.Errorf("no resolver"),
			},
			true,
			"unable to resolve: no resolver",
		},
		{
			"unregistered-name",
			args{
				errors.Errorf("unregistered name"),
			},
			true,
			"not found: unregistered name",
		},
		{
			"could-not-parse-address",
			args{
				errors.Errorf("could not parse address"),
			},
			true,
			"invalid name: could not parse address",
		},
		{
			"unknown",
			args{
				errors.Errorf("unknown error"),
			},
			true,
			"unknown error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wrapError(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("wrapError() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && (tt.wantErrMessage != "") && !assert.EqualError(err, tt.wantErrMessage) {
				t.Errorf("wrapError() errorMessage = %v, wantErrMessage %v", err, tt.wantErrMessage)
			}
			// if assert.
		})
	}
}
