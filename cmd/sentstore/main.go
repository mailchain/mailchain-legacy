package main

import (
	"log"

	"github.com/mailchain/mailchain/cmd/sentstore/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatalln(err)
	}
}
