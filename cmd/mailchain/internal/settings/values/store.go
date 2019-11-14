package values

//go:generate mockgen -source=store.go -package=valuestest -destination=./valuestest/store_mock.go

// Store provides the functions to retrieve values from configuration storage.
type Store interface {
	GetString(key string) string
	GetStringSlice(key string) []string
	GetInt(key string) int
	GetBool(key string) bool
	IsSet(key string) bool
	Set(key string, value interface{})
}
