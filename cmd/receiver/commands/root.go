package commands

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/mailchain/mailchain/cmd/receiver/handler"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "receiver",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := newDatabaseConnection(cmd)
			if err != nil {
				return errors.WithStack(err)
			}
			ts, err := pq.NewTransactionStore(db)
			if err != nil {
				return errors.WithStack(err)
			}

			r := mux.NewRouter()
			r.HandleFunc("/to", handler.HandleToRequest(ts)).Methods("GET")

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

	cmd.PersistentFlags().String("postgres-host", "localhost", "Postgres server host")
	cmd.PersistentFlags().String("postgres-sslmode", "disable", "Use SSL when connecting to Postgres")
	cmd.PersistentFlags().Int("postgres-port", 5432, "Postgres server port")

	cmd.PersistentFlags().String("postgres-user", "envelope", "Envelopes postgres database user")
	cmd.PersistentFlags().String("postgres-password", "", "Envelopes postgres database password")
	cmd.PersistentFlags().String("postgres-name", "envelope", "Envelopes postgres database name")

	return cmd
}
