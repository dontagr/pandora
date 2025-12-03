package interfaces

import (
	"github.com/dontagr/pandora/internal/server/store/models"
)

type (
	UserStore interface {
		GetUser(login string) (*models.User, error)
		SaveUser(login string, passwordHash string) error
	}
	SecretStore interface {
		GetSecret(userId int, kind string, label string) (*models.Secret, error)
		GetListSecret(userId int, kind string) (*models.SecretList, error)
		SyncSecret(userId int) (*models.SecretSync, error)
		DeleteSecret(userId int, kind string, label string) error
		SaveSecret(userId int, kind string, label string, data string, meta string, version int) (*models.Secret, error)
	}
)
