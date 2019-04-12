package stores

import (
	"errors"

	ldberr "github.com/syndtr/goleveldb/leveldb/errors"
)

var (
	errNotFound = errors.New("not found")
)

func IsNotFoundError(err error) bool {
	switch err {
	case errNotFound,
		ldberr.ErrNotFound:
		return true
	default:
		return false
	}
}
