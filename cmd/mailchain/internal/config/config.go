package config

import (
	"github.com/spf13/viper" // nolint: depguard
)

//go:generate mockgen -source=sent_store.go -package=configtest -destination=./configtest/sent_store_mock.go

type SentStoreSetter interface {
	Set(sentType string) error
}

type SentStore struct {
	viper         *viper.Viper
	requiredInput func(label string) (string, error)
}

type Keystore struct {
	viper                    *viper.Viper
	requiredInputWithDefault func(label string, defaultValue string) (string, error)
}
