package commands

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/nameservice/handler"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/nameservice/ens"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

func addENS(router *mux.Router) error {
	ethMainnet, err := ens.NewLookupService("https://mainnet.infura.io/v3/.....")
	if err != nil {
		return err
	}

	addItem(router, protocols.Ethereum, ethereum.Mainnet, ethMainnet)

	return nil
}

func addItem(router *mux.Router, network, protocol string, service nameservice.Lookup) {
	router.HandleFunc(fmt.Sprintf("/%s/%s/name", protocols.Ethereum, ethereum.Mainnet), handler.Forward(service, protocols.Ethereum, ethereum.Mainnet)).Methods("GET")
	router.HandleFunc(fmt.Sprintf("/%s/%s/address", protocols.Ethereum, ethereum.Mainnet), handler.Reverse(service, protocols.Ethereum, ethereum.Mainnet)).Methods("GET")
}

// {protocol}/{network}/name/?domain-name={domain-name}
// {protocol}/{network}/address?address={address}
func rootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use: "nameservice",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := mux.NewRouter()

			if err := addENS(r); err != nil {
				return err
			}

			n := negroni.New()
			n.UseHandler(r)
			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				return err
			}
			n.Run(fmt.Sprintf(":%d", port))
			return nil
		},
	}
	cmd.PersistentFlags().Int("port", 8080, "")

	return cmd, nil
}
