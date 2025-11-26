package store

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dontagr/pandora/internal/agent/service/kinds/factory"
	"github.com/dontagr/pandora/internal/agent/store/secret"
	"github.com/dontagr/pandora/internal/models"
)

type Service struct {
	kindFactory *factory.KindFactory
	secretStore *secret.Secret
}

func NewService(kindFactory *factory.KindFactory, secretStore *secret.Secret) *Service {
	return &Service{kindFactory: kindFactory, secretStore: secretStore}
}

func (us *Service) Encrypt(req *models.RequestStoreSave) error {
	kindService, err := us.kindFactory.GetKind(req.Kind)
	if err != nil {
		return fmt.Errorf("getKind: %v", err)
	}

	req.Data, err = kindService.Encrypt(req.Data)
	if err != nil {
		return fmt.Errorf("encrypt: %v", err)
	}

	return nil
}

func (us *Service) Decrypt(resp *models.ResponceStoreLoad, reveryData any) error {
	kindService, err := us.kindFactory.GetKind(resp.Kind)
	if err != nil {
		return fmt.Errorf("getKind: %v", err)
	}

	resp.ReveryData, err = kindService.Decrypt(reveryData)
	if err != nil {
		return fmt.Errorf("decrypt: %v", err)
	}

	return nil
}

func (us *Service) ResponceSync(body []byte) (*models.ResponceSync, error) {
	var storeLoad models.ResponceSync

	err := json.Unmarshal(body, &storeLoad)
	if err != nil {
		return nil, fmt.Errorf("error deserializing JSON: %v", err)
	}

	return &storeLoad, nil
}

func (us *Service) RecoveryStoreLoad(body []byte) (*models.ResponceStoreLoad, error) {
	var storeLoad models.ResponceStoreLoad

	err := json.Unmarshal(body, &storeLoad)
	if err != nil {
		return nil, fmt.Errorf("error deserializing JSON: %v", err)
	}

	kindService, err := us.kindFactory.GetKind(storeLoad.Kind)
	if err != nil {
		return nil, fmt.Errorf("getKind: %v", err)
	}

	reveryData, err := kindService.UnmarshalData(storeLoad.Data)
	if err != nil {
		return nil, fmt.Errorf("unmarshalData: %v", err)
	}

	err = us.Decrypt(&storeLoad, reveryData)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %v", err)
	}

	return &storeLoad, nil
}

func (us *Service) RecoveryStoreList(body []byte) (*models.ResponceStoreList, error) {
	var storeList models.ResponceStoreList

	err := json.Unmarshal(body, &storeList)
	if err != nil {
		return nil, fmt.Errorf("error deserializing JSON: %v", err)
	}

	return &storeList, nil
}

func (us *Service) GetRequestStoreLoad(cmd *cobra.Command) (*models.RequestStoreLoad, error) {
	kind, _ := cmd.Flags().GetString("type")
	label, _ := cmd.Flags().GetString("label")
	if kind == "" {
		return nil, fmt.Errorf("type is required")
	}
	kindService, err := us.kindFactory.GetKind(kind)
	if err != nil {
		return nil, err
	}
	err = kindService.ValidateLabel(label)
	if err != nil {
		return nil, err
	}

	return &models.RequestStoreLoad{Kind: kind, Label: label}, nil
}

func (us *Service) GetRequestStoreList(cmd *cobra.Command) (*models.RequestStoreList, error) {
	kind, _ := cmd.Flags().GetString("type")
	if kind == "" {
		return nil, fmt.Errorf("type is required")
	}
	_, err := us.kindFactory.GetKind(kind)
	if err != nil {
		return nil, err
	}

	return &models.RequestStoreList{Kind: kind}, nil
}

func (us *Service) GetRequestStoreSave(cmd *cobra.Command) (*models.RequestStoreSave, error) {
	kind, _ := cmd.Flags().GetString("type")
	data, _ := cmd.Flags().GetString("data")
	file, _ := cmd.Flags().GetString("file")
	label, _ := cmd.Flags().GetString("label")
	meta, _ := cmd.Flags().GetString("meta")
	if kind == "" {
		return nil, fmt.Errorf("type is required")
	}
	kindService, err := us.kindFactory.GetKind(kind)
	if err != nil {
		return nil, err
	}
	if kindService.IsDataRequired() {
		err = kindService.ValidateData(data)
		if err != nil {
			return nil, err
		}
	}
	err = kindService.ValidateLabel(label)
	if err != nil {
		return nil, err
	}
	err = kindService.ValidateMeta(meta)
	if err != nil {
		return nil, err
	}
	if kindService.IsFileRequired() {
		err = kindService.ValidateFile(file)
		if err != nil {
			return nil, err
		}
	} else {
		file = ""
	}

	sData, err := kindService.GetStruct(data, file)
	if err != nil {
		return nil, err
	}

	return &models.RequestStoreSave{Kind: kind, Data: sData, Label: label, Meta: meta}, nil
}

func (us *Service) GetRequestStoreSync(login string) (*models.RequestSyncList, error) {
	syncList, err := us.secretStore.GetSecretForSync(login)
	if err != nil {
		return nil, fmt.Errorf("GetSecretForSync: %v", err)
	}

	return syncList, nil
}

func (us *Service) UpdateLocalChanges(login string, kind string, label string, version int) error {
	err := us.secretStore.UpdateSecret(login, kind, label, version, 1)
	if err != nil {
		return fmt.Errorf("SaveSecret: %v", err)
	}

	return nil
}

func (us *Service) SaveLocalChanges(req *models.RequestStoreSave, synced bool, login string) error {
	data, err := req.GetData()
	if err != nil {
		return fmt.Errorf("GetData: %v", err)
	}

	err = us.secretStore.SaveSecret(login, req.Kind, req.Label, data, req.Meta, req.Version, boolToInt(synced))
	if err != nil {
		return fmt.Errorf("SaveSecret: %v", err)
	}

	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
