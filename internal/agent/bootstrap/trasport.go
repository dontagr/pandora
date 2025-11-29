package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/agent/service/transport"
)

// Transport предоставляет функции инициализации транспортного уровня приложения.
// Он использует библиотеку fx для предоставления зависимостей, связанных с транспортным слоем.
var Transport = fx.Options(
	// Инициализирует новый HTTPManager, который управляет HTTP соединениями и запросами.
	fx.Provide(transport.NewHTTPManager),
)
