package kindAuth

import (
	"fmt"
	"os"

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
		name: models.KindBinary,
		Kind: kindCommon.Kind{
			NeedFile: true,
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

func (m *Kind) GetStruct(_ string, file string) (any, error) {
	fileData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return &models.RequestKindBinary{Text: string(fileData)}, nil
}

func (m *Kind) Encrypt(data any) (any, error) {
	var err error
	req := data.(*models.RequestKindBinary)

	if req.Text != "" {
		req.Text, err = m.EncryptData(req.Text)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}
