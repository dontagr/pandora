package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/server/httpserver"
	"github.com/dontagr/pandora/internal/server/httpserver/routing"
)

var Server = fx.Options(
	fx.Provide(
		httpserver.NewServer,
	),
	fx.Invoke(
		func(*httpserver.HTTPServer) {},
		routing.InitRouting,
	),
)
