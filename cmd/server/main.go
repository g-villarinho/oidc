package main

import (
	"log"

	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server"
	"github.com/g-villarinho/oidc-server/pkg/injector"
)

func main() {
	container := server.InitializeContainer()

	app := injector.Resolve[*server.Server](container)

	if err := app.Start(); err != nil {
		log.Fatal(err.Error())
	}
}
