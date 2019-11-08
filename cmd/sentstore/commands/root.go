package commands

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/sentstore/handlers"
	"github.com/mailchain/mailchain/cmd/sentstore/storage"
	"github.com/mailchain/mailchain/stores"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

func rootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "sentstore",
		Short: "Mailchain sent store",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := cmd.Flags().GetString("base-url")
			if err != nil {
				return err
			}
			r := mux.NewRouter()
			region, err := cmd.Flags().GetString("aws-region")
			if err != nil {
				return err
			}
			bucket, err := cmd.Flags().GetString("aws-bucket")
			if err != nil {
				return err
			}

			store, err := storage.NewSentStore(
				region,
				bucket,
				"",
				"",
			)
			if err != nil {
				return err
			}
			r.HandleFunc("/", handlers.PostHandler(base, store, stores.SizeMegabyte*2)).Methods("POST")

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
	cmd.PersistentFlags().String("base-url", "", "")
	cmd.PersistentFlags().String("aws-region", "us-east-1", "")
	cmd.PersistentFlags().String("aws-bucket", "", "")
	cmd.PersistentFlags().Int("port", 8080, "")

	return cmd, nil
}
