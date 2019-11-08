package values //nolint:dupl

import (
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
)

//go:generate mockgen -source=int.go -package=valuestest -destination=./valuestest/int_mock.go

// Int interface to all Get, Set and exporting or `int` values.
type Int interface {
	Get() int
	Set(v int)
	Attribute() output.Attribute
}

// DefaultInt implementation of `Int` interface.
type DefaultInt struct {
	def     int
	setting string
	store   Store
}

// Get the value if set otherwise return default value.
func (d DefaultInt) Get() int {
	if d.store.IsSet(d.setting) {
		return d.store.GetInt(d.setting)
	}
	return d.def
}

// Set the value and store it in configuration store.
func (d DefaultInt) Set(v int) {
	d.store.Set(d.setting, v)
}

// Attribute creates representation of the value to be used when exporting the configuration.
func (d DefaultInt) Attribute() output.Attribute {
	dots := strings.Split(d.setting, ".")

	return output.Attribute{
		FullName:  dots[len(dots)-1],
		IsDefault: d.Get() == d.def,
		Value:     d.Get(),
	}
}

// NewDefaultInt create the `Int` value.
func NewDefaultInt(defVal int, store Store, setting string) Int {
	return DefaultInt{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
