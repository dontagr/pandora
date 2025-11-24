package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/agent/service/transport"
)

var Transport = fx.Options(
	fx.Provide(transport.NewHTTPManager),
)
