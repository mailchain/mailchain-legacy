package ens

import (
	"fmt"

	"github.com/mailchain/mailchain/internal/nameservice"
	"github.com/pkg/errors"
)

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	switch fmt.Sprintf("%v", errors.Cause(err)) {
	case fmt.Sprintf("%v", errors.Cause(errors.Errorf("no resolver"))),
		fmt.Sprintf("%v", errors.Cause(errors.Errorf("No resolution"))):
		return errors.WithMessage(err, nameservice.ErrUnableToResolve.Error())

	case fmt.Sprintf("%v", errors.Cause(errors.Errorf("unregistered name"))):
		return errors.WithMessage(err, nameservice.ErrNotFound.Error())

	case fmt.Sprintf("%v", errors.Cause(errors.Errorf("could not parse address"))):
		// address related to not being able to part ens address not ethereum address
		return errors.WithMessage(err, nameservice.ErrInvalidName.Error())

	default:
		return err
	}
}
