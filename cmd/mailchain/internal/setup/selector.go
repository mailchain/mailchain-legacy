package setup

//go:generate mockgen -source=selector.go -package=setuptest -destination=./setuptest/selector_mock.go

type ChainNetworkExistingSelector interface {
	Select(chain, network, receiver string) (string, error)
}
