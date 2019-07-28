package commands

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/nameservice/handler"
	"github.com/mailchain/mailchain/internal/nameservice"
	"github.com/mailchain/mailchain/internal/nameservice/ens"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

func config() (map[string]nameservice.Lookup, error) {
	ethMainnet, err := ens.NewLookupService("https://mainnet.infura.io/v3/.....")
	if err != nil {
		return nil, err
	}
	return map[string]nameservice.Lookup{
		"ethereum/mainnet": ethMainnet,
		// ...
	}, nil
}

// {protocol}/{network}/name/?domain-name={domain-name}
// {protocol}/{network}/address?address={address}
func rootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use: "nameservice",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := mux.NewRouter()
			for k, v := range config() {
				println(fmt.Sprintf("%s/name", k))
				r.HandleFunc(fmt.Sprintf("/%s/name", k), handler.Forward(v)).Methods("GET")
				r.HandleFunc(fmt.Sprintf("/%s/address", k), handler.Reverse(v)).Methods("GET")
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
