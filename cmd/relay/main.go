package main

import (
	"log"

	"github.com/mailchain/mailchain/cmd/relay/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatalln(err)
	}
}
