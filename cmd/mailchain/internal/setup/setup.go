package setup

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/spf13/viper" // nolint: depguard
)

type Keystore struct {
	setter             config.KeystoreSetter
	viper              *viper.Viper
	selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
}
type Network struct {
	receiverSelector     ChainNetworkExistingSelector
	senderSelector       ChainNetworkExistingSelector
	pubKeyFinderSelector ChainNetworkExistingSelector
	selectItem           func(label string, items []string) (string, error)
}
type PubKeyFinder struct {
	setter             config.PubKeyFinderSetter
	viper              *viper.Viper
	selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
}

type Receiver struct {
	setter             config.ReceiverSetter
	viper              *viper.Viper
	selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
}

type Sender struct {
	setter             config.SenderSetter
	viper              *viper.Viper
	selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
}

type SentStorage struct {
	setter             config.SentStoreSetter
	viper              *viper.Viper
	selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
}
