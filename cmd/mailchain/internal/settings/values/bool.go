package values //nolint:dupl

import (
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
)

//go:generate mockgen -source=bool.go -package=valuestest -destination=./valuestest/bool_mock.go

// Bool interface to all Get, Set and exporting or `bool` values.
type Bool interface {
	Get() bool
	Set(v bool)
	Attribute() output.Attribute
}

// DefaultBool implementation of `Bool` interface
type DefaultBool struct {
	def     bool
	setting string
	store   Store
}

// Get the value if set otherwise return default value.
func (d DefaultBool) Get() bool {
	if d.store.IsSet(d.setting) {
		return d.store.GetBool(d.setting)
	}

	return d.def
}

// Set the value and store it in configuration store.
func (d DefaultBool) Set(v bool) {
	d.store.Set(d.setting, v)
}

// Attribute creates representation of the value to be used when exporting the configuration.
func (d DefaultBool) Attribute() output.Attribute {
	dots := strings.Split(d.setting, ".")

	return output.Attribute{
		FullName:  dots[len(dots)-1],
		IsDefault: d.Get() == d.def,
		Value:     d.Get(),
	}
}

// NewDefaultBool create the `Bool` value.
func NewDefaultBool(defVal bool, store Store, setting string) Bool {
	return DefaultBool{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
