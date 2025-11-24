package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/server/service/interfaces"
	"github.com/dontagr/pandora/internal/server/store/user"
)

var Store = fx.Options(
	fx.Provide(
		fx.Annotate(
			user.NewUser,
			fx.As(new(interfaces.UserStore)),
		),
	),
	fx.Invoke(
		func(interfaces.UserStore) {},
	),
)
