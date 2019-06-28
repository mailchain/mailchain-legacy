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
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/handlers"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus" // nolint:depguard
	"github.com/spf13/cobra"
	"github.com/spf13/viper" // nolint:depguard
	"github.com/urfave/negroni"
)

func CreateRouter(s *settings.Base, cmd *cobra.Command) (http.Handler, error) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/spec.json", handlers.GetSpec()).Methods("GET")
	api.HandleFunc("/docs", handlers.GetDocs()).Methods("GET")
	receivers, err := s.Ethereum.GetReceivers(s.Receivers)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not configure receivers")
	}
	pubKeyFinders, err := s.Ethereum.GetPublicKeyFinders(s.PublicKeyFinders)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not configure receivers")
	}

	senders, err := s.Ethereum.GetSenders(s.Senders)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not configure senders")
	}

	sentStorage, err := s.SentStore.Produce()
	if err != nil {
		return nil, errors.WithMessage(err, "Could not config store")
	}
	mailboxStore, err := s.MailboxState.Produce()
	if err != nil {
		return nil, errors.WithMessage(err, "Could not config mailbox store")
	}
	keystore, err := s.Keystore.Produce()
	if err != nil {
		return nil, errors.WithMessage(err, "could not create `keystore`")
	}
	passphrase, err := config.Passphrase(cmd)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get `passphrase`")
	}
	// TODO: currently this only does scrypt need flag + config etc
	deriveKeyOptions := multi.OptionsBuilders{
		Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase(passphrase)},
	}

	api.HandleFunc("/addresses", handlers.GetAddresses(keystore)).Methods("GET")

	api.HandleFunc("/ethereum/{network}/address/{address:[-0-9a-zA-Z]+}/public-key", handlers.GetPublicKey(pubKeyFinders)).Methods("GET")
	api.HandleFunc(
		"/ethereum/{network}/address/{address:[-0-9a-zA-Z]+}/messages",
		handlers.GetMessages(mailboxStore, receivers, keystore, deriveKeyOptions)).Methods("GET")
	api.HandleFunc("/ethereum/{network}/messages/send",
		handlers.SendMessage(sentStorage, senders, keystore, deriveKeyOptions)).Methods("POST")
	api.HandleFunc("/messages/{message_id}/read", handlers.GetRead(mailboxStore)).Methods("GET")
	api.HandleFunc("/messages/{message_id}/read", handlers.PutRead(mailboxStore)).Methods("PUT")
	api.HandleFunc("/messages/{message_id}/read", handlers.DeleteRead(mailboxStore)).Methods("DELETE")

	_ = r.Walk(gorillaWalkFn)
	return r, nil
}

func SetupFlags(cmd *cobra.Command) error {
	cmd.Flags().Int("port", defaults.Port, "Port to run server on")
	cmd.Flags().Bool("cors-disabled", defaults.CORSDisabled, "Disable CORS on the server")

	if err := viper.BindPFlag("server.port", cmd.Flags().Lookup("port")); err != nil {
		return err
	}
	if err := viper.BindPFlag("server.cors.disabled", cmd.Flags().Lookup("cors-disabled")); err != nil {
		return err
	}
	// if err := viper.BindPFlag("server.cors.allowed-origins", cmd.Flags().Lookup("cors-allowed-origins")); err != nil {
	// 	return err
	// }

	cmd.PersistentFlags().String("passphrase", "", "Passphrase to encrypt/decrypt key with")
	return nil
}

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
