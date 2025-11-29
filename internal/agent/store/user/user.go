// Package user предоставляет функциональность для управления данными пользователя,
// включая хранение и извлечение информации о пользователе из базы данных SQLite.
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

// SQL-запросы для создания и управления таблицей пользователей в базе данных.
const (
	// createUserTable создает таблицу пользователей, если она еще не существует.
	createUserTable = `CREATE TABLE IF NOT EXISTS user ("id" INTEGER NOT NULL, "login" TEXT NOT NULL, "token" TEXT NOT NULL, CONSTRAINT "uq_id" UNIQUE ("id"));`
	// insertUserSQL вставляет или обновляет запись пользователя в таблице.
	insertUserSQL = `INSERT INTO user ("id", "login", "token") VALUES (1, ?, ?) ON CONFLICT("id") DO UPDATE SET login=excluded.login, token=excluded.token;`
	// selectUserSQL выбирает запись пользователя из таблицы.
	selectUserSQL = `SELECT "login", "token" FROM user WHERE "id"=1 LIMIT 1;`
)

// User предоставляет методы для управления данными пользователя в базе данных.
type User struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

// NewUser создает новый экземпляр User и инициализирует таблицу пользователя в базе данных.
// Он также регистрирует хуки запуска для управления схемой базы данных и загрузки пользовательских данных.
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

// addShema создаёт таблицу пользователя в базе данных, если она еще не существует.
func (u *User) addShema() error {
	if _, err := u.db.Exec(createUserTable); err != nil {
		return err
	}

	return nil
}

// SaveUserAuth сохраняет информацию об авторизации пользователя, включая логин и токен, в базе данных.
func (u *User) SaveUserAuth(login string, token string) error {
	_, err := u.db.Exec(insertUserSQL, login, token)

	return err
}

// GetUserAuth извлекает данные об авторизации пользователя из базы данных.
// Он возвращает структуру User с логином и токеном пользователя.
func (u *User) GetUserAuth() (*models.User, error) {
	rows := u.db.QueryRow(selectUserSQL)

	var us models.User
	err := rows.Scan(&us.Login, &us.Token)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetUserAuth Scan: %w", err)
	}

	return &us, nil
}
