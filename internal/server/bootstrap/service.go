package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/server/service/jwt"
	"github.com/dontagr/pandora/internal/server/service/user"
)

var Service = fx.Options(
	fx.Provide(
		user.NewUserService,
		jwt.NewJWTService,
	),
)
