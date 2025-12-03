package bootstrap

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"

	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/agent/config"
	"github.com/dontagr/pandora/internal/agent/store/secret"
	"github.com/dontagr/pandora/internal/agent/store/user"
)

// SqlLite объединяет предоставление зависимостей для SQL-базы данных,
// пользовательского и секретного хранилищ, а также их инициализацию с использованием fx.
var SqlLite = fx.Options(
	fx.Provide(
		newSqlLite,
		user.NewUser,
		secret.NewSecret,
	),
	fx.Invoke(
		func(*user.User) {},
		func(*secret.Secret) {},
	),
)

// newSqlLite создает и настраивает подключение к базе данных SQLite.
// Она также добавляет хуки для открытия и закрытия соединения при старте и остановке приложения.
func newSqlLite(cfg *config.Config, lc fx.Lifecycle) (*sql.DB, error) {
	// Открывает новое подключение к базе данных по указанному пути
	db, err := sql.Open("sqlite", cfg.DataBase.Path)
	if err != nil {
		return nil, err
	}

	// Добавляет хуки старта и остановки для управления жизненным циклом подключения к базе данных
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
