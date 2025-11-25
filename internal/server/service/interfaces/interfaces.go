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
		GetSecret(userId int, kind string) (*models.Secret, error)
		SaveSecret(userId int, kind string, label string, data string) error
	}
)
