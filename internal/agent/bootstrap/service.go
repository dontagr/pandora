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

// Service организует предоставление и регистрацию всех необходимых сервисов
// и их зависимостей в приложении, используя библиотеку fx.
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

// NewCManager создает и возвращает новый экземпляр CManager для управления криптографическими ключами.
// Он инициализирует публичный и приватный ключи, используя данные из конфига.
func NewCManager(cfg *config.Config) (*crypro.CManager, error) {
	cmanager := crypro.CManager{}
	// Инициализация публичного ключа
	err := cmanager.InitPublicKey(cfg.Crypto.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("NewHasher: %w", err)
	}
	// Инициализация приватного ключа
	err = cmanager.InitPrivateKey(cfg.Crypto.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("NewHasher: %w", err)
	}

	return &cmanager, nil
}
