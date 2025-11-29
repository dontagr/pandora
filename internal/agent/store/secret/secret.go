// Package secret предоставляет функциональность для управления секретами пользователей,
// включая сохранение, загрузку и синхронизацию секретов в базе данных.
package secret

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/pandora/internal/models"
)

// SQL-запросы для создания и управления таблицей секретов в базе данных.
const (
	// createSecretTable создает таблицу секретов, если она еще не существует.
	createSecretTable = `CREATE TABLE IF NOT EXISTS secret (
    	"login" TEXT NOT NULL, 
		"kind" TEXT NOT NULL,
		"label" TEXT NOT NULL,
		"data" TEXT NOT NULL,
		"meta" TEXT NOT NULL,
		"version" INTEGER NOT NULL,
		"sync" INTEGER NOT NULL,
		CONSTRAINT "uq_login_kind_label" UNIQUE ("login", "kind", "label")
    );`
	// insertSecretSQL вставляет или обновляет запись секретов в таблице.
	insertSecretSQL = `INSERT INTO secret ("login", "kind", "label", "data", "meta", "version", "sync") VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT("login", "kind", "label") DO UPDATE SET "data"=excluded.data, "meta"=excluded.meta, "sync"=excluded.sync, "version"=excluded.version;`
	// updateSecretSQL обновляет версию и статус синхронизации имеющегося секрета в таблице.
	updateSecretSQL = `UPDATE secret SET "version" = ?, "sync" = ? WHERE "login" = ? AND "kind" = ? AND "label" = ?;`
	// selectSecretForSyncSQL выбирает все секреты, которые нуждаются в синхронизации.
	selectSecretForSyncSQL = `SELECT "kind", "label", "data", "meta", "version" FROM secret WHERE "login"=? AND "sync"=0;`
	// selectVersionSecretSQL выбирает версию секрета из таблицы.
	selectVersionSecretSQL = `SELECT "version" FROM secret WHERE "login"=? AND "kind" = ? AND "label" = ?;`
	// selectSecretSQL выбирает все данные определенного секрета.
	selectSecretSQL = `SELECT "data", "meta", "version" FROM secret WHERE "login"=? AND "kind" = ? AND "label" = ?;`
)

// Secret предоставляет методы для управления секретами в базе данных.
type Secret struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

// NewSecret создает новый экземпляр Secret и инициализирует таблицу секретов в базе данных.
// Он также регистрирует хук для инициализации схемы базы данных при запуске приложения.
func NewSecret(log *zap.SugaredLogger, db *sql.DB, lc fx.Lifecycle) *Secret {
	secret := Secret{
		db:  db,
		log: log,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := secret.addShema()
			if err != nil {
				return err
			}

			return nil
		},
	})

	return &secret
}

// addShema создает таблицу секретов в базе данных, если она еще не существует.
func (u *Secret) addShema() error {
	if _, err := u.db.Exec(createSecretTable); err != nil {
		return err
	}

	return nil
}

// LoadSecret загружает секрет из базы данных по предоставленным логину, типу и метке.
func (u *Secret) LoadSecret(login string, kind string, label string) (*models.RequestStoreSave, error) {
	rows := u.db.QueryRow(selectSecretSQL, login, kind, label)

	store := models.RequestStoreSave{Kind: kind, Label: label}
	err := rows.Scan(&store.Data, &store.Meta, &store.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("selectVersion Scan: %w", err)
	}

	return &store, nil
}

// GetSecretVersion получает текущую версию секрета из базы данных.
func (u *Secret) GetSecretVersion(login string, kind string, label string) (int, error) {
	rows := u.db.QueryRow(selectVersionSecretSQL, login, kind, label)

	var version int
	err := rows.Scan(&version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("selectVersion Scan: %w", err)
	}

	return version, nil
}

// UpdateSecret обновляет версию и статус синхронизации секрета в базе данных.
func (u *Secret) UpdateSecret(login string, kind string, label string, version int, sync int) error {
	_, err := u.db.Exec(updateSecretSQL, version, sync, login, kind, label)

	return err
}

// SaveSecret сохраняет новый секрет или обновляет существующий секрет в базе данных.
func (u *Secret) SaveSecret(login string, kind string, label string, data string, meta string, version int, sync int) error {
	_, err := u.db.Exec(insertSecretSQL, login, kind, label, data, meta, version, sync)

	return err
}

// GetSecretForSync получает список секретов, которые нужно синхронизировать с удалённым сервером.
func (u *Secret) GetSecretForSync(login string) (*models.RequestSyncList, error) {
	rows, err := u.db.Query(selectSecretForSyncSQL, login)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	result := &models.RequestSyncList{List: make([]*models.RequestStoreSave, 0)}
	for rows.Next() {
		node := models.RequestStoreSave{}

		err = rows.Scan(&node.Kind, &node.Label, &node.Data, &node.Meta, &node.Version)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		node.Version = node.Version + 1

		result.List = append(result.List, &node)
	}

	return result, nil
}
