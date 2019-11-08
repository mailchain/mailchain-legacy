package values //nolint:dupl

import (
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
)

//go:generate mockgen -source=string_slice.go -package=valuestest -destination=./valuestest/string_slice_mock.go

// StringSlice interface to all Get, Set and exporting of `[]string` values.
type StringSlice interface {
	Get() []string
	Set(v []string)
	Attribute() output.Attribute
}

// DefaultStringSlice implementation of `StringSlice` interface.
type DefaultStringSlice struct {
	def     []string
	setting string
	store   Store
}

// Get the value if set otherwise return default value.
func (d DefaultStringSlice) Get() []string {
	if d.store.IsSet(d.setting) {
		return d.store.GetStringSlice(d.setting)
	}

	return d.def
}

// Set the value and store it in configuration store.
func (d DefaultStringSlice) Set(v []string) {
	d.store.Set(d.setting, v)
}

// Attribute creates representation of the value to be used when exporting the configuration.
func (d DefaultStringSlice) Attribute() output.Attribute {
	dots := strings.Split(d.setting, ".")

	return output.Attribute{
		FullName:  dots[len(dots)-1],
		IsDefault: strings.Join(d.Get(), "-") == strings.Join(d.def, "-"),
		Value:     d.Get(),
	}
}

// NewDefaultStringSlice create the `StringSlice` value.
func NewDefaultStringSlice(defVal []string, store Store, setting string) StringSlice {
	return DefaultStringSlice{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
