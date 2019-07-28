package nameservice

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrUnableToResolve = errors.New("unable to resolve")
	ErrNotFound        = errors.New("not found")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidAddress  = errors.New("invalid address")
)

func IsNoResolverError(err error) bool {
	message := fmt.Sprintf("%v", err)
	if strings.HasPrefix(message, ErrUnableToResolve.Error()) {
		return true
	}
	return false
}

func IsNotFoundError(err error) bool {
	message := fmt.Sprintf("%v", err)
	if strings.HasPrefix(message, ErrNotFound.Error()) {
		return true
	}

	return false
}

func IsInvalidNameError(err error) bool {
	message := fmt.Sprintf("%v", err)
	if strings.HasPrefix(message, ErrInvalidName.Error()) {
		return true
	}

	return false
}

func IsInvalidAddressError(err error) bool {
	message := fmt.Sprintf("%v", err)
	if strings.HasPrefix(message, ErrInvalidAddress.Error()) {
		return true
	}
	return false
}
