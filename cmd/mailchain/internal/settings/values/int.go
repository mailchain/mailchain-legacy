// nolint:dupl
package values

//go:generate mockgen -source=int.go -package=valuestest -destination=./valuestest/int_mock.go
type Int interface {
	Get() int
	Set(v int)
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

func NewDefaultInt(defVal int, store Store, setting string) Int {
	return DefaultInt{
		def:     defVal,
		setting: setting,
		store:   store,
	}
}
