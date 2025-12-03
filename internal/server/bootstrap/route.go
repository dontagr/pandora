package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/server/service/handler"
)

var Route = fx.Options(
	fx.Provide(handler.NewHandler),
	fx.Invoke(func(handler *handler.Handler) {}),
)
