// nolint:dupl
package values

//go:generate mockgen -source=string_slice.go -package=valuestest -destination=./valuestest/string_slice_mock.go
type StringSlice interface {
	Get() []string
	Set(v []string)
}

type DefaultStringSlice struct {
	def     []string
	setting string
	store   Store
}

func (d DefaultStringSlice) Get() []string {
	if d.store.IsSet(d.setting) {
		return d.store.GetStringSlice(d.setting)
	}
	return d.def
}

func (d DefaultStringSlice) Set(v []string) {
	d.store.Set(d.setting, v)
}

func NewDefaultStringSlice(defVal []string, store Store, setting string) StringSlice {
	return DefaultStringSlice{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
