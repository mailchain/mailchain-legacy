// nolint:dupl
package values

//go:generate mockgen -source=string.go -package=valuestest -destination=./valuestest/string_mock.go
type String interface {
	Get() string
	Set(v string)
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

func NewDefaultString(defVal string, store Store, setting string) String {
	return DefaultString{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
