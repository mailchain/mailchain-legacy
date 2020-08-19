package cachestore

import (
	"os"
	"time"

	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/hash"

	"github.com/spf13/afero"
)

type CacheStore struct {
	cache afero.Fs
}

func NewCacheStore(cacheTimeout time.Duration, basePath string) *CacheStore {
	base := afero.NewBasePathFs(afero.NewOsFs(), basePath)
	layer := afero.NewMemMapFs()
	cache := afero.NewCacheOnReadFs(base, layer, cacheTimeout)
	return &CacheStore{cache: cache}
}

func (f *CacheStore) GetMessage(location string) ([]byte, error) {
	key := f.getCacheKey(location)
	return afero.ReadFile(f.cache, key)
}

func (f *CacheStore) SetMessage(location string, msg []byte) error {
	key := f.getCacheKey(location)
	if _, err := f.cache.Stat(key); os.IsNotExist(err) {
		_, err := f.cache.Create(key)
		if err != nil {
			return err
		}
	}
	return afero.WriteFile(f.cache, key, msg, os.ModePerm)
}

func (f *CacheStore) getCacheKey(location string) string {
	locHash, _ := hash.Create(hash.SHA3256, []byte(location))
	return encoding.EncodeHex(locHash)
}

//cleanUp is used for cleaning up test files
func (f *CacheStore) cleanUp(location string) error {
	key := f.getCacheKey(location)
	return f.cache.Remove(key)
}
