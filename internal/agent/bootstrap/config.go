package bootstrap

import (
	"fmt"

	"go.uber.org/fx"

	configInternal "github.com/dontagr/pandora/internal/agent/config"
	"github.com/dontagr/pandora/pkg/config"
)

var Config = fx.Options(
	fx.Provide(newConfig),
)

func newConfig() (*configInternal.Config, error) {
	configInt := configInternal.NewConfig()

	cnf := &config.Config{
		Data:             configInt,
		DefaultFilePaths: []string{"../../../configs", "./configs"},
		DefaultFileNames: []string{"agent.config.json"},
	}

	cnf.ReadFromFile()
	err := cnf.ReadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to read from env: %w", err)
	}

	err = cnf.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return configInt, nil
}
