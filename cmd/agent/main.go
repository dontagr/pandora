package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/agent/bootstrap"
)

func main() {
	app := fx.New(CreateApp())

	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}
}

func CreateApp() fx.Option {
	return fx.Options(
		bootstrap.Config,
		bootstrap.Logger,
		bootstrap.SqlLite,
		bootstrap.Transport,
		bootstrap.Service,
		bootstrap.CLI,
	)
}
