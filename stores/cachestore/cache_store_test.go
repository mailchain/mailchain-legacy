package cachestore

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	type args struct {
		location string
		cache    func() *CacheStore
	}
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
		args
	}{
		{
			name: "success",
			want: []byte("input file is small"),
			args: args{
				location: "location/location",
				cache: func() *CacheStore {
					cacheStore := NewCacheStore(10*time.Second, "../testdata")
					err := cacheStore.SetMessage("location/location", []byte("input file is small"))
					if err != nil {
						t.Fatal(err)
					}
					return cacheStore
				},
			},
		},
		{
			name: "cache miss",
			want: nil,
			args: args{
				location: "location/location",
				cache: func() *CacheStore {
					return NewCacheStore(10*time.Second, "../testdata")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.args.cache()
			defer func() {
				if err := cache.cleanUp(tt.args.location); err != nil {
					fmt.Println(err)
				}
			}()
			got, err := cache.GetMessage(tt.args.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ReadMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
