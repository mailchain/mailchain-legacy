// nolint:dupl
package values

import (
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
)

//go:generate mockgen -source=int.go -package=valuestest -destination=./valuestest/int_mock.go
type Int interface {
	Get() int
	Set(v int)
	Attribute() output.Attribute
}

type DefaultInt struct {
	def     int
	setting string
	store   Store
}

func (d DefaultInt) Get() int {
	if d.store.IsSet(d.setting) {
		return d.store.GetInt(d.setting)
	}
	return d.def
}

func (d DefaultInt) Set(v int) {
	d.store.Set(d.setting, v)
}

func (d DefaultInt) Attribute() output.Attribute {
	dots := strings.Split(d.setting, ".")

	return output.Attribute{
		Name:      dots[len(dots)-1],
		IsDefault: d.Get() == d.def,
		Value:     d.Get(),
	}
}

func NewDefaultInt(defVal int, store Store, setting string) Int {
	return DefaultInt{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
