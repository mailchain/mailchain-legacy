package commands

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/relay/handler"
	"github.com/mailchain/mailchain/cmd/relay/relayer"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

func config() map[string]relayer.RelayFunc {
	return map[string]relayer.RelayFunc{
		"ethereum/ropsten": relayer.ChangeURL("URL"),
		"ethereum/mainnet": relayer.ChangeURL("URL"),
		// ...
	}
}

func rootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use: "relay",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := mux.NewRouter()
			r.PathPrefix("/").HandlerFunc(handler.HandleRequest(config()))

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
