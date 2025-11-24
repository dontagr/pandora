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
			ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
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
