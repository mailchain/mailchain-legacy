// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/handlers"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus" //nolint:depguard
	"github.com/spf13/cobra"
	"github.com/spf13/viper" //nolint:depguard
	"github.com/ttacon/chalk"
	"github.com/urfave/negroni"
)

// CreateRouter configure a router with all api resources.
func CreateRouter(s *settings.Root, cmd *cobra.Command) (http.Handler, error) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/spec.json", handlers.GetSpec()).Methods("GET")
	api.HandleFunc("/docs", handlers.GetDocs()).Methods("GET")

	config, err := produceConfig(s)
	if err != nil {
		return nil, errors.WithMessage(err, "could not config http server")
	}

	cmdPassphrase, _ := cmd.Flags().GetString("passphrase")
	passphrase, err := prompts.Secret(cmdPassphrase,
		fmt.Sprint(chalk.Yellow, "Note: To derive a storage key passphrase is required. The passphrase must be secure and not guessable."),
		"Passphrase",
		false,
		true,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get `passphrase`")
	}
	// TODO: currently this only does scrypt need flag + config etc
	deriveKeyOptions := multi.OptionsBuilders{
		Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase(passphrase)},
	}

	api.HandleFunc("/addresses", handlers.GetAddresses(config.keystore)).Methods("GET")

	api.HandleFunc("/messages", handlers.GetMessages(config.mailboxStateStore, config.receivers, config.keystore, deriveKeyOptions)).Methods("GET") //nolint:lll
	api.HandleFunc("/messages", handlers.SendMessage(config.sentStore, config.senders, config.keystore, deriveKeyOptions)).Methods("POST")

	api.HandleFunc("/messages/{message_id}/read", handlers.GetRead(config.mailboxStateStore)).Methods("GET")
	api.HandleFunc("/messages/{message_id}/read", handlers.PutRead(config.mailboxStateStore)).Methods("PUT")
	api.HandleFunc("/messages/{message_id}/read", handlers.DeleteRead(config.mailboxStateStore)).Methods("DELETE")

	api.HandleFunc("/protocols", handlers.GetProtocols(s)).Methods("GET")

	api.HandleFunc("/public-key", handlers.GetPublicKey(config.publicKeyFinders)).Methods("GET")

	api.HandleFunc("/nameservice/address/{address}/resolve", handlers.GetResolveAddress(config.addressResolvers)).Methods("GET")
	api.HandleFunc("/nameservice/name/{domain-name}/resolve", handlers.GetResolveName(config.domainResolvers)).Methods("GET")

	api.HandleFunc("/version", handlers.GetVersion()).Methods("GET")

	_ = r.Walk(gorillaWalkFn)
	return r, nil
}

// SetupFlags created default flags and bind's the values to configuration settings.
func SetupFlags(cmd *cobra.Command) error {
	cmd.Flags().Int("port", defaults.Port, "Port to run server on")
	cmd.Flags().Bool("cors-disabled", defaults.CORSDisabled, "Disable CORS on the server")

	if err := viper.BindPFlag("server.port", cmd.Flags().Lookup("port")); err != nil {
		return err
	}
	if err := viper.BindPFlag("server.cors.disabled", cmd.Flags().Lookup("cors-disabled")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("passphrase", "", "Passphrase to encrypt/decrypt key with")
	return nil
}

// CreateNegroni returns a server, with CORS and endpoints configured
func CreateNegroni(config *settings.Server, router http.Handler) *negroni.Negroni {
	n := negroni.New()

	if !config.CORS.Disabled.Get() {
		n.Use(cors.New(cors.Options{
			AllowedOrigins: config.CORS.AllowedOrigins.Get(),
			AllowedHeaders: []string{"Authorization", "Content-Type"},
			AllowedMethods: []string{"GET", "PUT", "DELETE", "POST", "HEAD", "PATCH"},
			MaxAge:         86400,
		}))
	}

	n.UseHandler(router)

	return n
}

func gorillaWalkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	methods, _ := route.GetMethods()
	for _, method := range methods {
		path, _ := route.GetPathTemplate()
		log.Infof("Serving %s : %s", method, path)
	}
	return nil
}
