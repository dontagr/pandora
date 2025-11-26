package secret

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/pandora/internal/models"
)

const (
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
	insertSecretSQL        = `INSERT INTO secret ("login", "kind", "label", "data", "meta", "version", "sync") VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT("login", "kind", "label") DO UPDATE SET "data"=excluded.data, "meta"=excluded.meta, "sync"=excluded.sync, "version"=excluded.version;`
	updateSecretSQL        = `UPDATE secret SET "version" = ?, "sync" = ? WHERE "login" = ? AND "kind" = ? AND "label" = ?;`
	selectSecretForSyncSQL = `SELECT "kind", "label", "data", "meta", "version" FROM secret WHERE "login"=? AND "sync"=0;`
)

type Secret struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

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

func (u *Secret) addShema() error {
	if _, err := u.db.Exec(createSecretTable); err != nil {
		return err
	}

	return nil
}

func (u *Secret) UpdateSecret(login string, kind string, label string, version int, sync int) error {
	_, err := u.db.Exec(updateSecretSQL, version, sync, login, kind, label)

	return err
}

func (u *Secret) SaveSecret(login string, kind string, label string, data string, meta string, version int, sync int) error {
	_, err := u.db.Exec(insertSecretSQL, login, kind, label, data, meta, version, sync)

	return err
}

func (u *Secret) GetSecretForSync(login string) (*models.RequestSyncList, error) {
	rows, err := u.db.Query(selectSecretForSyncSQL, login)
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}
	defer rows.Close()

	result := &models.RequestSyncList{List: make([]models.SyncNode, 0)}
	for rows.Next() {
		node := models.SyncNode{}

		err = rows.Scan(&node.Kind, &node.Label, &node.Data, &node.Meta, &node.Version)
		if err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		result.List = append(result.List, node)
	}

	return result, nil
}
