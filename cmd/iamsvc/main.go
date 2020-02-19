package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/tierklinik-dobersberg/identity-server/cmd/iamsvc/app"
)

var (
	defaultPort = "8080"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cmd := app.NewIAMCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
