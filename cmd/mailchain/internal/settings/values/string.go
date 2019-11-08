package values //nolint:dupl

import (
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
)

//go:generate mockgen -source=string.go -package=valuestest -destination=./valuestest/string_mock.go

// String interface to all Get, Set and exporting or `string` values.
type String interface {
	Get() string
	Set(v string)
	Attribute() output.Attribute
}

// DefaultString implementation of `String` interface.
type DefaultString struct {
	def     string
	setting string
	store   Store
}

// Get the value if set otherwise return default value.
func (d DefaultString) Get() string {
	if d.store.IsSet(d.setting) {
		return d.store.GetString(d.setting)
	}

	return d.def
}

// Set the value and store it in configuration store.
func (d DefaultString) Set(v string) {
	d.store.Set(d.setting, v)
}

// Attribute creates representation of the value to be used when exporting the configuration.
func (d DefaultString) Attribute() output.Attribute {
	dots := strings.Split(d.setting, ".")

	return output.Attribute{
		FullName:  dots[len(dots)-1],
		IsDefault: d.Get() == d.def,
		Value:     d.Get(),
	}
}

// NewDefaultString create the `String` value.
func NewDefaultString(defVal string, store Store, setting string) String {
	return DefaultString{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
