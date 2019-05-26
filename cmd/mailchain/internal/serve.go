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

package internal

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/mailchain/config"
	"github.com/mailchain/mailchain/cmd/mailchain/config/defaults"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/addresses"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/ethereum/address/messages"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/ethereum/address/publickey"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/ethereum/messages/send"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/messages/read"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/spec"
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/scrypt"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus" // nolint:depguard
	"github.com/spf13/cobra"
	"github.com/spf13/viper" // nolint:depguard
	"github.com/urfave/negroni"
)

func CreateRouter(cmd *cobra.Command) (http.Handler, error) {
	r := mux.NewRouter()
	r.HandleFunc("/api/spec.json", spec.Get()).Methods("GET")
	r.HandleFunc("/api/docs", spec.DocsGet()).Methods("GET")
	vpr := viper.GetViper()
	receivers, err := config.GetReceivers(vpr)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not configure receivers")
	}
	pubKeyFinders, err := config.GetPublicKeyFinders(vpr)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not configure receivers")
	}
	senders, err := config.GetSenders(vpr)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not configure senders")
	}

	sentStorage, err := config.GetSentStorage()
	if err != nil {
		return nil, errors.WithMessage(err, "Could not config store")
	}
	mailboxStore, err := config.GetStateStore(vpr)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not config mailbox store")
	}
	keystore, err := config.GetKeystore()
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
	r.HandleFunc("/api/addresses", addresses.GetAddresses(keystore)).Methods("GET")
	r.HandleFunc("/api/ethereum/{network}/address/{address:[-0-9a-zA-Z]+}/public-key", publickey.GetPublicKey(pubKeyFinders)).Methods("GET")
	r.HandleFunc(
		"/api/ethereum/{network}/address/{address:[-0-9a-zA-Z]+}/messages",
		messages.GetMessages(mailboxStore, receivers, keystore, deriveKeyOptions)).Methods("GET")
	r.HandleFunc("/api/ethereum/{network}/messages/send", send.SendMessage(sentStorage, senders, keystore, deriveKeyOptions)).Methods("POST")
	r.HandleFunc("/api/messages/{message_id}/read", read.GetRead(mailboxStore)).Methods("GET")
	r.HandleFunc("/api/messages/{message_id}/read", read.PutRead(mailboxStore)).Methods("PUT")
	r.HandleFunc("/api/messages/{message_id}/read", read.DeleteRead(mailboxStore)).Methods("DELETE")

	_ = r.Walk(gorillaWalkFn)
	return r, nil
}

func SetupFlags(cmd *cobra.Command) error {
	cmd.Flags().Int("port", defaults.Port, "Port to run server on")
	cmd.Flags().Bool("cors-disabled", defaults.CORSDisabled, "Disable CORS on the server")
	cmd.Flags().String("cors-allowed-origins", defaults.CORSAllowedOrigins, "Allowed origins for CORS")

	if err := viper.BindPFlag("server.port", cmd.Flags().Lookup("port")); err != nil {
		return err
	}
	if err := viper.BindPFlag("server.cors.disabled", cmd.Flags().Lookup("cors-disabled")); err != nil {
		return err
	}
	if err := viper.BindPFlag("server.cors.allowed-origins", cmd.Flags().Lookup("cors-allowed-origins")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("passphrase", "", "Passphrase to encrypt/decrypt key with")
	return nil
}

func CreateNegroni(router http.Handler) *negroni.Negroni {
	n := negroni.New()
	if !viper.GetBool("server.cors.disabled") {
		n.Use(cors.New(cors.Options{
			AllowedOrigins: []string{viper.GetString("server.cors.allowed-origins")},
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
