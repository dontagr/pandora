package secret

import (
	"github.com/dontagr/pandora/internal/server/service/interfaces"
	"github.com/dontagr/pandora/internal/server/service/jwt"
	"github.com/dontagr/pandora/internal/server/store/models"
)

type Service struct {
	store      interfaces.SecretStore
	jwtService *jwt.JWTService
}

func NewSecretService(store interfaces.SecretStore, jwtService *jwt.JWTService) *Service {
	return &Service{store: store, jwtService: jwtService}
}

func (u *Service) GetSecret(userId int, kind string) (*models.Secret, error) {
	return u.store.GetSecret(userId, kind)
}

func (u *Service) SaveSecret(userId int, kind string, label string, data string) error {
	return u.store.SaveSecret(userId, kind, label, data)
}
