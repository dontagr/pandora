package secret

import (
	output "github.com/dontagr/pandora/internal/models"
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

func (u *Service) GetResponceStoreList(secret *models.SecretList) *output.ResponceStoreList {
	resList := &output.ResponceStoreList{
		List: make([]output.StoreLite, len(secret.List)),
	}

	for i, secretLite := range secret.List {
		resList.List[i] = output.StoreLite{
			Label:   secretLite.Label,
			DT:      secretLite.DT,
			Version: secretLite.Version,
		}
	}

	return resList
}

func (u *Service) GetResponceStoreSync(secret *models.SecretSync) *output.ResponceSyncList {
	result := output.ResponceSyncList{List: make([]*output.ResponceStoreLoad, len(secret.List))}

	for i, row := range secret.List {
		result.List[i] = &output.ResponceStoreLoad{
			Kind:    row.Kind,
			Label:   row.Label,
			Version: row.Version,
			Data:    row.Data,
			Meta:    row.Meta,
			DT:      row.DT,
		}
	}

	return &result
}

func (u *Service) GetResponceStoreLoad(secret *models.Secret) *output.ResponceStoreLoad {
	return &output.ResponceStoreLoad{
		Kind:    secret.Kind,
		Label:   secret.Label,
		Version: secret.Version,
		Data:    secret.Data,
		Meta:    secret.Meta,
		DT:      secret.DT,
	}
}

func (u *Service) SyncSecret(userId int) (*models.SecretSync, error) {
	return u.store.SyncSecret(userId)
}

func (u *Service) GetListSecret(userId int, kind string) (*models.SecretList, error) {
	return u.store.GetListSecret(userId, kind)
}

func (u *Service) GetSecret(userId int, kind string, label string) (*models.Secret, error) {
	return u.store.GetSecret(userId, kind, label)
}

func (u *Service) SaveSecret(userId int, kind string, label string, data string, meta string, version int) (*models.Secret, error) {
	return u.store.SaveSecret(userId, kind, label, data, meta, version)
}

func (u *Service) DeleteSecret(userId int, kind string, label string) error {
	return u.store.DeleteSecret(userId, kind, label)
}
