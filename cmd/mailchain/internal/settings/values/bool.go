// nolint:dupl
package values

//go:generate mockgen -source=bool.go -package=valuestest -destination=./valuestest/bool_mock.go
type Bool interface {
	Get() bool
	Set(v bool)
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

func NewDefaultBool(defVal bool, store Store, setting string) Bool {
	return DefaultBool{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
