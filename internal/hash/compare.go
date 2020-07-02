package hash

import (
	"bytes"

	"github.com/pkg/errors"
)

func CompareContentsToHash(data, hash []byte) error {
	kind, digest, err := parse(hash)
	if err != nil {
		return errors.Wrap(err, "compare hash: parse failed")
	}

	h, err := Create(kind, data)
	if err != nil {
		return errors.Wrap(err, "compare hash: create failed")
	}

	digestContents, err := GetDigest(kind, h)
	if err != nil {
		return errors.Wrap(err, "compare hash: get digest failed")
	}

	if !bytes.Equal(digest, digestContents) {
		return errors.Errorf("compare hash: hashes do not match")
	}

	return nil
}
