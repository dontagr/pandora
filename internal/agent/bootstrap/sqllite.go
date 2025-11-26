package bootstrap

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"

	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/agent/config"
	"github.com/dontagr/pandora/internal/agent/store/user"
)

var SqlLite = fx.Options(
	fx.Provide(
		newSqlLite,
		user.NewUser,
	),
	fx.Invoke(func(*user.User) {}),
)

func newSqlLite(cfg *config.Config, lc fx.Lifecycle) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.DataBase.Path)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()
			if err = db.PingContext(ctx); err != nil {
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			db.Close()
			return nil
		},
	})

	return db, nil
}
