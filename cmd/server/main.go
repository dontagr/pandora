package main

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/server/bootstrap"
)

func main() {
	app := fx.New(CreateApp())

	//startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	//defer cancel()
	//if err := app.Start(startCtx); err != nil {
	//	log.Fatal(err)
	//}

	app.Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		bootstrap.Server,
		bootstrap.Config,
		bootstrap.Logger,
		bootstrap.Postgres,
		bootstrap.Store,
		bootstrap.Route,
		bootstrap.Service,
	)
}
