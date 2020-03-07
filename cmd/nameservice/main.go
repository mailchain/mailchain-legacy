package main

import (
	"log"

	"github.com/mailchain/mailchain/cmd/nameservice/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
