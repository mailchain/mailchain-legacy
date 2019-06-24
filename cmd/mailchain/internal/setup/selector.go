package setup

//go:generate mockgen -source=selector.go -package=setuptest -destination=./setuptest/selector_mock.go

type ChainNetworkExistingSelector interface {
	Select(chain, network, selectorType string) (selectedItem string, err error)
}

type SimpleSelector interface {
	Select(selectorType string) (string, error)
}

type KeystoreSelector interface {
	Select(keystoreType, keystorePath string) (string, error)
}
