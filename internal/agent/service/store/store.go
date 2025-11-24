package store

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dontagr/pandora/internal/agent/service/kinds/factory"
	"github.com/dontagr/pandora/internal/models"
)

type Service struct {
	kindFactory *factory.KindFactory
}

func NewService(kindFactory *factory.KindFactory) *Service {
	return &Service{kindFactory: kindFactory}
}

func (us *Service) Encrypt(req *models.RequestStoreSave) (*models.RequestStoreSave, error) {
	kindService, err := us.kindFactory.GetKind(req.Kind)
	if err != nil {
		return nil, fmt.Errorf("GetKind: ", err)
	}

	req.Data, err = kindService.Encrypt(req.Data)
	if err != nil {
		return nil, fmt.Errorf("Encrypt: ", err)
	}

	return req, nil
}

func (us *Service) GetRequestStoreSave(cmd *cobra.Command) (*models.RequestStoreSave, error) {
	kind, _ := cmd.Flags().GetString("type")
	data, _ := cmd.Flags().GetString("data")
	file, _ := cmd.Flags().GetString("file")
	label, _ := cmd.Flags().GetString("label")
	if kind == "" {
		return nil, fmt.Errorf("type is required")
	}
	kindService, err := us.kindFactory.GetKind(kind)
	if err != nil {
		return nil, err
	}
	if kindService.IsDataRequired() {
		err := kindService.ValidateData(data)
		if err != nil {
			return nil, err
		}
	}
	err = kindService.ValidateLabel(label)
	if err != nil {
		return nil, err
	}
	if kindService.IsFileRequired() {
		err := kindService.ValidateFile(file)
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

	return &models.RequestStoreSave{Kind: kind, Data: sData, Label: label}, nil
}
