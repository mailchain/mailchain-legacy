package stores

//go:generate mockgen -source=cache.go -package=storestest -destination=./storestest/cache_mock.go

// Cache allows you to cache the message for a certain of time.
type Cache interface {
	GetMessage(location string) ([]byte, error)
	SetMessage(location string, msg []byte) error
}
