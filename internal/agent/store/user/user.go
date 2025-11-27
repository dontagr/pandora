package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/pandora/internal/agent/service/transport"
	"github.com/dontagr/pandora/internal/agent/store/models"
)

const (
	createUserTable = `CREATE TABLE IF NOT EXISTS user ("id" INTEGER NOT NULL, "login" TEXT NOT NULL, "token" TEXT NOT NULL, CONSTRAINT "uq_id" UNIQUE ("id"));`
	insertUserSQL   = `INSERT INTO user ("id", "login", "token") VALUES (1, ?, ?) ON CONFLICT("id") DO UPDATE SET login=excluded.login, token=excluded.token;`
	selectUserSQL   = `SELECT "login", "token" FROM user WHERE "id"=1 LIMIT 1;`
)

type User struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

func NewUser(log *zap.SugaredLogger, db *sql.DB, lc fx.Lifecycle, client *transport.HTTPManager) *User {
	user := User{
		db:  db,
		log: log,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := user.addShema()
			if err != nil {
				return err
			}

			usr, err := user.GetUserAuth()
			if err != nil {
				return err
			}
			client.User = usr

			return nil
		},
	})

	return &user
}

func (u *User) addShema() error {
	if _, err := u.db.Exec(createUserTable); err != nil {
		return err
	}

	return nil
}

func (u *User) SaveUserAuth(login string, token string) error {
	_, err := u.db.Exec(insertUserSQL, login, token)

	return err
}

func (u *User) GetUserAuth() (*models.User, error) {
	rows := u.db.QueryRow(selectUserSQL)

	var us models.User
	err := rows.Scan(&us.Login, &us.Token)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetUserAuth Scan: %w", err)
	}

	return &us, nil
}
