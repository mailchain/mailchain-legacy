// nolint:dupl
package values

import (
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
)

//go:generate mockgen -source=bool.go -package=valuestest -destination=./valuestest/bool_mock.go
type Bool interface {
	Get() bool
	Set(v bool)
	Attribute() output.Attribute
}

type DefaultBool struct {
	def     bool
	setting string
	store   Store
}

func (d DefaultBool) Get() bool {
	if d.store.IsSet(d.setting) {
		return d.store.GetBool(d.setting)
	}

	return d.def
}

func (d DefaultBool) Set(v bool) {
	d.store.Set(d.setting, v)
}

func (d DefaultBool) Attribute() output.Attribute {
	dots := strings.Split(d.setting, ".")

	return output.Attribute{
		FullName:  dots[len(dots)-1],
		IsDefault: d.Get() == d.def,
		Value:     d.Get(),
	}
}

func NewDefaultBool(defVal bool, store Store, setting string) Bool {
	return DefaultBool{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
