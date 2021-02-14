package algorand

const (
	// Mainnet network name.
	Mainnet = "mainnet"

	// Betanet network name.
	Betanet = "betanet"

	// Testnet network name.
	Testnet = "testnet"
)

// Networks supported by substrate package.
func Networks() []string {
	return []string{Mainnet, Betanet, Testnet}
}
