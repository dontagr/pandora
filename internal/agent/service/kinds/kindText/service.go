package kindAuth

import (
	"encoding/json"
	"fmt"

	"github.com/dontagr/pandora/internal/agent/service/kinds/factory"
	"github.com/dontagr/pandora/internal/agent/service/kinds/kindCommon"
	"github.com/dontagr/pandora/internal/models"
	crypro "github.com/dontagr/pandora/pkg/crypto"
)

type Kind struct {
	kindCommon.Kind
	name string
}

func RegisterKind(mf *factory.KindFactory, hasher *crypro.CManager) {
	mf.SetKind(&Kind{
		name: models.KindText,
		Kind: kindCommon.Kind{
			NeedData: true,
			Hasher:   hasher,
		},
	})
}

func (m *Kind) GetName() string {
	return m.name
}

func (m *Kind) ValidateData(data string) error {
	err := m.CommonValidateData(data)
	if err != nil {
		return err
	}

	return nil
}

func (m *Kind) GetStruct(data string, _ string) (any, error) {
	var req models.RequestKindText

	req.Text = data

	return &req, nil
}

func (m *Kind) Encrypt(data any) (any, error) {
	var err error
	req := data.(*models.RequestKindText)

	if req.Text != "" {
		req.Text, err = m.EncryptData(req.Text)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (m *Kind) Decrypt(data any) (map[string]string, error) {
	var err error
	req := data.(*models.RequestKindText)

	if req.Text != "" {
		req.Text, err = m.DecryptData(req.Text)
		if err != nil {
			return nil, err
		}
	}

	return m.StructToMap(*req), nil
}

func (m *Kind) UnmarshalData(data string) (any, error) {
	var reqData models.RequestKindText

	err := json.Unmarshal([]byte(data), &reqData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal data: %w", err)
	}

	return &reqData, nil
}
