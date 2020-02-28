package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/tierklinik-dobersberg/identity-server/cmd/iamcli/cmds"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := cmds.RootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
