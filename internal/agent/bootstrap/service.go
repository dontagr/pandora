package bootstrap

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/agent/config"
	"github.com/dontagr/pandora/internal/agent/service/kinds/factory"
	kindAuth "github.com/dontagr/pandora/internal/agent/service/kinds/kindAuth"
	kindBinary "github.com/dontagr/pandora/internal/agent/service/kinds/kindBinary"
	kindCard "github.com/dontagr/pandora/internal/agent/service/kinds/kindCard"
	kindText "github.com/dontagr/pandora/internal/agent/service/kinds/kindText"
	"github.com/dontagr/pandora/internal/agent/service/store"
	"github.com/dontagr/pandora/internal/agent/service/user"
	crypro "github.com/dontagr/pandora/pkg/crypto"
)

var Service = fx.Options(
	fx.Provide(
		user.NewService,
		store.NewService,
		factory.NewKindFactory,
		NewCManager,
	),
	fx.Invoke(
		kindAuth.RegisterKind,
		kindCard.RegisterKind,
		kindBinary.RegisterKind,
		kindText.RegisterKind,
	),
)

func NewCManager(cfg *config.Config) (*crypro.CManager, error) {
	cmanager := crypro.CManager{}
	err := cmanager.InitPublicKey(cfg.Crypto.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("NewHasher: %v", err)
	}
	err = cmanager.InitPrivateKey(cfg.Crypto.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("NewHasher: %v", err)
	}

	return &cmanager, nil
}
