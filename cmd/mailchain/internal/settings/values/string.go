// nolint:dupl
package values

import (
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
)

//go:generate mockgen -source=string.go -package=valuestest -destination=./valuestest/string_mock.go
type String interface {
	Get() string
	Set(v string)
	Attribute() output.Attribute
}

type DefaultString struct {
	def     string
	setting string
	store   Store
}

func (d DefaultString) Get() string {
	if d.store.IsSet(d.setting) {
		return d.store.GetString(d.setting)
	}
	return d.def
}

func (d DefaultString) Set(v string) {
	d.store.Set(d.setting, v)
}

func (d DefaultString) Attribute() output.Attribute {
	dots := strings.Split(d.setting, ".")

	return output.Attribute{
		Name:      dots[len(dots)-1],
		IsDefault: d.Get() == d.def,
		Value:     d.Get(),
	}
}

func NewDefaultString(defVal string, store Store, setting string) String {
	return DefaultString{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
