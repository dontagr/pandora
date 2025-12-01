package bootstrap

import (
	"fmt"

	"go.uber.org/fx"

	configInternal "github.com/dontagr/pandora/internal/agent/config"
	"github.com/dontagr/pandora/pkg/config"
)

// Config определяет набор настроек для инициализации конфигурации приложения.
// Используется библиотекой fx для предоставления зависимостей конфигурации.
var Config = fx.Options(
	fx.Provide(newConfig),
)

// newConfig создает и инициализирует конфигурацию приложения.
// Она загружает настройки из файлов, окружения и проверяет их корректность.
// Возвращает указатель на структуру конфигурации или ошибку, если что-то пошло не так.
func newConfig() (*configInternal.Config, error) {
	configInt := configInternal.NewConfig()

	cnf := &config.Config{
		Data:             configInt,
		DefaultFilePaths: []string{"../../../configs", "./configs"},
		DefaultFileNames: []string{"agent.config.json"},
	}

	// Загружает конфигурацию из файла
	err := cnf.ReadFromFile()
	if err != nil {
		return nil, fmt.Errorf("ReadFromFile: %v", err)
	}

	// Читает переменные окружения для переопределения или дополнения
	err = cnf.ReadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to read from env: %w", err)
	}

	// Проверяет корректность загруженной конфигурации
	err = cnf.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return configInt, nil
}
