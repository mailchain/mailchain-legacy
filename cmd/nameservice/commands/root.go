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
	ethMainnet, err := ens.NewLookupService("https://eth-mainnet.alchemyapi.io/jsonrpc/jjeXif1lwcfY_J_fpvbds55mZWuoXkFD")
	if err != nil {
		return err
	}

	addItem(router, protocols.Ethereum, ethereum.Mainnet, ethMainnet)

	return nil
}

func addItem(router *mux.Router, protocol, network string, service nameservice.Lookup) {
	router.HandleFunc(fmt.Sprintf("/%s/%s/name", protocol, network), handler.Forward(service, protocol, network)).Methods("GET")
	router.HandleFunc(fmt.Sprintf("/%s/%s/address", protocol, network), handler.Reverse(service, protocol, network)).Methods("GET")
}

// {protocol}/{network}/name/?domain-name={domain-name}
// {protocol}/{network}/address?address={address}
func rootCmd() *cobra.Command {
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

	return cmd
}
