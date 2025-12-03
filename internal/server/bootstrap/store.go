package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/server/service/interfaces"
	"github.com/dontagr/pandora/internal/server/store/secret"
	"github.com/dontagr/pandora/internal/server/store/user"
)

var Store = fx.Options(
	fx.Provide(
		fx.Annotate(
			user.NewUser,
			fx.As(new(interfaces.UserStore)),
		),
		fx.Annotate(
			secret.NewSecret,
			fx.As(new(interfaces.SecretStore)),
		),
	),
	fx.Invoke(
		func(interfaces.UserStore) {},
	),
)
