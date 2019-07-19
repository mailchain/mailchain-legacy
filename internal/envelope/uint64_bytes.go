package envelope

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

type UInt64Bytes []byte

func NewUInt64Bytes(i uint64, b []byte) UInt64Bytes {
	loc := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(loc, i)

	buf := make([]byte, len(loc[:n])+1+len(b))
	buf[0] = byte(n)
	copy(buf[1:], loc[:n])
	copy(buf[1+n:], b)

	return buf
}

func (u UInt64Bytes) UInt64() (uint64, error) {
	i, _, err := parseUInt64Bytes(u)
	return i, err
}

func (u UInt64Bytes) Bytes() ([]byte, error) {
	_, b, err := parseUInt64Bytes(u)
	return b, err
}

func (u UInt64Bytes) Values() (i uint64, b []byte, err error) {
	return parseUInt64Bytes(u)
}

func parseUInt64Bytes(buf []byte) (i uint64, b []byte, err error) {
	if len(buf) == 0 {
		return 0, []byte{}, errors.Errorf("\"buf\" must not be empty")
	}
	bufLen := int(buf[0])
	if len(buf) < bufLen+1 {
		return 0, []byte{}, errors.Errorf("\"buf\" is too short to be valid")
	}
	intPortion := buf[1 : 1+bufLen]
	v, n := binary.Uvarint(intPortion)
	if n != len(intPortion) {
		return 0, []byte{}, errors.Errorf("uint64 did not consume all data")
	}
	return v, buf[1+bufLen:], nil
}
