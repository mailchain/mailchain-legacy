package scrypt

import "golang.org/x/crypto/scrypt"

func DeriveKey(o []DeriveOptionsBuilder) ([]byte, error) {
	opts := &DeriveOpts{}
	apply(opts, o)
	return scrypt.Key([]byte(opts.Passphrase), opts.Salt, opts.N, opts.R, opts.P, opts.Len)
}

func CreateOptions(o []DeriveOptionsBuilder) *DeriveOpts {
	opts := &DeriveOpts{}
	apply(opts, o)
	return opts
}

func apply(o *DeriveOpts, opts []DeriveOptionsBuilder) {
	for _, f := range opts {
		f(o)
	}
}
